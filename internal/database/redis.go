package database

import (
	"context"
	"fmt"
	"ganjineh-auth/internal/config"
	"strconv"
	"strings"
	"time"

	"github.com/google/wire"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

type ServiceRedisInterface interface {
	Health() map[string]string
	Client() *redis.Client
}

type ServiceRedisStruct struct {
	Db *redis.Client
}

func NewRedisConnection(conf *config.StructConfig) *ServiceRedisStruct {
	
	host := conf.REDIS_HOST + ":" + strconv.Itoa(conf.REDIS_PORT)
	rdb := redis.NewClient(&redis.Options{
		Addr: host,
		// Username: conf.REDIS_USER,
		// Password: conf.REDIS_PASS,
		DB: 0,
	})
	fmt.Print("redis run")
	return &ServiceRedisStruct{Db: rdb} 
}

var RedisSet = wire.NewSet(
	NewRedisConnection,
	wire.Bind(new(ServiceRedisInterface),new(*ServiceRedisStruct)),
)

var _ ServiceRedisInterface =(*ServiceRedisStruct)(nil)

// Health returns the health status and statistics of the Redis server.
func (s *ServiceRedisStruct) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Default is now 5s
	defer cancel()

	stats := make(map[string]string)

	// Check Redis health and populate the stats map
	stats = s.checkRedisHealth(ctx, stats)

	return stats
}

func (s *ServiceRedisStruct) Client() *redis.Client {
	return s.Db
}

// checkRedisHealth checks the health of the Redis server and adds the relevant statistics to the stats map.
func (s *ServiceRedisStruct) checkRedisHealth(ctx context.Context, stats map[string]string) map[string]string {
	// Ping the Redis server to check its availability.
	pong, err := s.Db.Ping(ctx).Result()
	if err != nil {
		stats["redis_status"] = "down"
		stats["redis_message"] = fmt.Sprintf("Redis connection failed: %v", err)
		return stats
	}

	// Redis is up
	stats["redis_status"] = "up"
	stats["redis_message"] = "It's healthy"
	stats["redis_ping_response"] = pong

	// Try to get Redis info, but don't fail if it errors
	info, err := s.Db.Info(ctx).Result()
	if err != nil {
		stats["redis_message"] = fmt.Sprintf("Connected but info failed: %v", err)
		return stats
	}

	// Continue with the rest of your health check...
	redisInfo := parseRedisInfo(info)
	// ... rest of the function
	return s.evaluateRedisStats(redisInfo, stats)
}

// evaluateRedisStats evaluates the Redis server statistics and updates the stats map with relevant messages.
func (s *ServiceRedisStruct) evaluateRedisStats(redisInfo, stats map[string]string) map[string]string {
	poolSize := s.Db.Options().PoolSize
	poolStats := s.Db.PoolStats()
	connectedClients, _ := strconv.Atoi(redisInfo["connected_clients"])
	highConnectionThreshold := int(float64(poolSize) * 0.8)

	// Check if the number of connected clients is high.
	if connectedClients > highConnectionThreshold {
		stats["redis_message"] = "Redis has a high number of connected clients"
	}

	// Check if the number of stale connections exceeds a threshold.
	minStaleConnectionsThreshold := 500
	if int(poolStats.StaleConns) > minStaleConnectionsThreshold {
		stats["redis_message"] = fmt.Sprintf("Redis has %d stale connections.", poolStats.StaleConns)
	}

	// Check if Redis is using a significant amount of memory.
	usedMemory, _ := strconv.ParseInt(redisInfo["used_memory"], 10, 64)
	maxMemory, _ := strconv.ParseInt(redisInfo["maxmemory"], 10, 64)
	if maxMemory > 0 {
		usedMemoryPercentage := float64(usedMemory) / float64(maxMemory) * 100
		if usedMemoryPercentage >= 90 {
			stats["redis_message"] = "Redis is using a significant amount of memory"
		}
	}

	// Check if Redis has been recently restarted.
	uptimeInSeconds, _ := strconv.ParseInt(redisInfo["uptime_in_seconds"], 10, 64)
	if uptimeInSeconds < 3600 {
		stats["redis_message"] = "Redis has been recently restarted"
	}

	// Check if the number of idle connections is high.
	idleConns := int(poolStats.IdleConns)
	highIdleConnectionThreshold := int(float64(poolSize) * 0.7)
	if idleConns > highIdleConnectionThreshold {
		stats["redis_message"] = "Redis has a high number of idle connections"
	}

	// Check if the connection pool utilization is high.
	poolUtilization := float64(poolStats.TotalConns-poolStats.IdleConns) / float64(poolSize) * 100
	highPoolUtilizationThreshold := 90.0
	if poolUtilization > highPoolUtilizationThreshold {
		stats["redis_message"] = "Redis connection pool utilization is high"
	}

	return stats
}

// parseRedisInfo parses the Redis info response and returns a map of key-value pairs.
func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result
}
