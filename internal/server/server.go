package server

import (
	"ganjineh-auth/internal/database"
	"ganjineh-auth/internal/routes"

	"github.com/gofiber/fiber/v2"

	// "ganjineh-auth/internal/database/redi"
	"ganjineh-auth/internal/repositories"
)

type FiberServer struct{
    *fiber.App
    pdb         database.ServicePostgresInterface
    rdb         database.ServiceRedisInterface  
    repos       *repositories.Container
}

func New( // ✅ Change to New to match wire.go
    app *fiber.App,
    pdb database.ServicePostgresInterface,
    rdb database.ServiceRedisInterface,
    repos *repositories.Container,
    routeContainer *routes.RouteContainer,
) *FiberServer {
    server := &FiberServer{
        App: app, // Use injected app, don't create new one
        pdb: pdb,
        rdb: rdb,
        repos: repos,
    }
    server.RegisterFiberRoutes(routeContainer)
    return server
}