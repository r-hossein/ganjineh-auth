package server

import (
	"context"
	"fmt"
	"ganjineh-auth/internal/database"
	"ganjineh-auth/internal/routes"
	"ganjineh-auth/internal/services"
	"ganjineh-auth/pkg"
	"time"

	"github.com/gofiber/fiber/v2"

	"ganjineh-auth/internal/repositories"
)

type FiberServer struct{
    *fiber.App
    pdb         database.ServicePostgresInterface
    rdb         database.ServiceRedisInterface  
    repos       *repositories.Container
    routeContainer *routes.RouteContainer
    startupService   services.StartupServiceInterface
}

func New( // ✅ Change to New to match wire.go
    pdb database.ServicePostgresInterface,
    rdb database.ServiceRedisInterface,
    repos *repositories.Container,
    routeContainer *routes.RouteContainer,
    errorHandler *pkg.ErrorHandler,
    startupService services.StartupServiceInterface,
) *FiberServer {
    app := fiber.New(fiber.Config{
        ServerHeader: "ganjineh-auth",
        AppName:      "ganjineh-auth",
        ErrorHandler: errorHandler.FiberErrorHandler,
    })
    server := &FiberServer{
        App: app, // Use injected app, don't create new one
        pdb: pdb,
        rdb: rdb,
        repos: repos,
        routeContainer: routeContainer,
        startupService: startupService,
    }
    
    server.initializeServer()
    
    server.RegisterFiberRoutes(routeContainer)
    return server
}

func (s *FiberServer) initializeServer() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := s.startupService.Initialize(ctx); err != nil {
        fmt.Printf("Warning: Failed to initialize roles: %v", err)
        // Don't panic, but log the error. The server can still start.
    } else {
        fmt.Println("Roles initialized successfully")
    }
}