package server

import (
	"ganjineh-auth/internal/database"
	"ganjineh-auth/internal/routes"
	"ganjineh-auth/pkg"

	"github.com/gofiber/fiber/v2"

	// "ganjineh-auth/internal/database/redi"
	"ganjineh-auth/internal/repositories"
)

type FiberServer struct{
    *fiber.App
    pdb         database.ServicePostgresInterface
    rdb         database.ServiceRedisInterface  
    repos       *repositories.Container
    routeContainer *routes.RouteContainer
}

func New( // ✅ Change to New to match wire.go
    pdb database.ServicePostgresInterface,
    rdb database.ServiceRedisInterface,
    repos *repositories.Container,
    routeContainer *routes.RouteContainer,
    errorHandler *pkg.ErrorHandler,
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
    }
    server.RegisterFiberRoutes(routeContainer)
    return server
}