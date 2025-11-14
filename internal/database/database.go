package database

import (
	"context"
	"fmt"
	"ganjineh-auth/internal/repositories/db"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type ServiceP interface {
	// Health returns a map of health status information.
	Health() map[string]string

	// Close terminates the database connection.
	Close() error
	
	// GetQueries returns the queries instance for database operations
	GetQueries() *db.Queries
	
	// GetPool returns the connection pool (برای موارد پیشرفته)
	GetPool() *pgxpool.Pool
}

type serviceP struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *serviceP
)

func NewP() ServiceP {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", 
		username, password, host, port, database, schema)

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

	// Create queries instance
	queries := db.New(pool)

	dbInstance = &serviceP{
		pool:    pool,
		queries: queries,
	}
	
	log.Printf("Connected to PostgreSQL database: %s", database)
	return dbInstance
}

// GetQueries returns the queries instance for database operations
func (s *serviceP) GetQueries() *db.Queries {
	return s.queries
}

// GetPool returns the connection pool
func (s *serviceP) GetPool() *pgxpool.Pool {
	return s.pool
}

// Health checks the health of the database connection
func (s *serviceP) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.pool.Ping(ctx)
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
	poolStats := s.pool.Stat()
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
func (s *serviceP) Close() error {
	log.Printf("Disconnecting from database: %s", database)
	s.pool.Close()
	return nil
}