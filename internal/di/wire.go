//go:build wireinject
// +build wireinject

package di

import (
	"ganjineh-auth/internal/config"
	"ganjineh-auth/internal/database"
	"ganjineh-auth/internal/handlers"
	"ganjineh-auth/internal/repositories"
	"ganjineh-auth/internal/routes"
	"ganjineh-auth/internal/server"
	"ganjineh-auth/internal/services"
	"ganjineh-auth/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

func InitializeUserHandler() (*server.FiberServer, error) {
    wire.Build(
        provideFiberApp,

		config.LoadConfig,
        utils.JwtPkgSet,
        utils.OTPPkgSet,
        utils.ValidatorSet,

        database.PostgreSQLSet,
        database.RedisSet,
        
        repositories.ContainerSet,

        services.AuthServiceSet,
        services.OTPServiceSet,

        handlers.AuthHandlerSet,
        
        server.ProvideGraphQLHandler,
        
        routes.RouteContainerSet,

        server.New,
    )
    return &server.FiberServer{}, nil
}

func provideFiberApp() *fiber.App {
    return fiber.New(fiber.Config{
        ServerHeader: "ganjineh-auth",
        AppName:      "ganjineh-auth",
    })
}