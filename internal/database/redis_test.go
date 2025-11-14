package database

import (
	"context"
	"log"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func mustStartRedisContainerC() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	dbContainer, err := redis.Run(
		context.Background(),
		"docker.io/redis:7.2.4",
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "6379/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	addressc = dbHost
	port = dbPort.Port()
	database = "0"

	return dbContainer.Terminate, err
}

func RTestMain(m *testing.M) {
	teardown, err := mustStartRedisContainerC()
	if err != nil {
		log.Fatalf("could not start redis container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown redis container: %v", err)
	}
}

func TestNewC(t *testing.T) {
	srv := NewR()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealthC(t *testing.T) {
	srv := NewR()

	stats := srv.Health()

	if stats["redis_status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["redis_status"])
	}

	if _, ok := stats["redis_version"]; !ok {
		t.Fatalf("expected redis_version to be present, got %v", stats["redis_version"])
	}
}
