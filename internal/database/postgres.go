package database

import (
	"context"
	"fmt"
	"ganjineh-auth/internal/config"
	"github.com/google/wire"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type ServicePostgresInterface interface {
	// Health returns a map of health status information.
	Health() map[string]string

	// Close terminates the database connection.
	Close() error
	
}

type ServicePostgresStruct struct {
	Pool    *pgxpool.Pool
}


func NewPostgresConnction(con *config.StructConfig) *ServicePostgresStruct {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", 
		con.BLUEPRINT_DB_USERNAME, con.BLUEPRINT_DB_PASSWORD, con.BLUEPRINT_DB_HOST, con.BLUEPRINT_DB_PORT, con.BLUEPRINT_DB_DATABASE, con.BLUEPRINT_DB_SCHEMA)

	// Create connection pool config
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("Failed to parse connection string: ", err)
	}

	// Configure connection pool
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 30 * time.Second
	config.HealthCheckPeriod = 1 * time.Minute

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Failed to create connection pool: ", err)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Failed to ping database: ", err)
	}
	
	log.Printf("Connected to PostgreSQL database")
	return &ServicePostgresStruct{
		Pool: pool,
	}
}



var PostgreSQLSet = wire.NewSet(
	NewPostgresConnction,
	wire.Bind(new(ServicePostgresInterface),new(*ServicePostgresStruct)),
)

var _ ServicePostgresInterface =(*ServicePostgresStruct)(nil)

// Health checks the health of the database connection
func (s *ServicePostgresStruct) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.Pool.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("db down: %v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get pool statistics
	poolStats := s.Pool.Stat()
	stats["total_connections"] = strconv.Itoa(int(poolStats.TotalConns()))
	stats["idle_connections"] = strconv.Itoa(int(poolStats.IdleConns()))
	stats["acquired_connections"] = strconv.Itoa(int(poolStats.AcquiredConns()))
	stats["max_connections"] = strconv.Itoa(int(poolStats.MaxConns()))
	// stats["empty_acquires"] = strconv.FormatUint(poolStats.EmptyAcquireCount(), 10)
	// stats["canceled_acquires"] = strconv.FormatUint(poolStats.CanceledAcquireCount(), 10)

	// Evaluate stats to provide a health message
	if poolStats.AcquiredConns() > 20 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if poolStats.EmptyAcquireCount() > 1000 {
		stats["message"] = "High number of empty acquires, consider increasing MaxConns."
	}

	return stats
}

// Close closes the database connection pool
func (s *ServicePostgresStruct) Close() error {
	log.Printf("Disconnecting from database")
	s.Pool.Close()
	return nil
}