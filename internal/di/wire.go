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
	"ganjineh-auth/internal/middleware"
	"ganjineh-auth/pkg"
	"github.com/google/wire"
)

func InitializeUserHandler() (*server.FiberServer, error) {
    wire.Build(
        pkg.NewErrorHandler,

		config.LoadConfig,
        utils.JwtPkgSet,
        utils.OTPPkgSet,
        utils.ValidatorSet,

        middleware.MiddlewareSet,
        
        database.PostgreSQLSet,
        database.RedisSet,
        
        repositories.ContainerSet,

        services.AuthServiceSet,
        services.OTPServiceSet,
        services.RewriteRoleServiceSet,
        services.StartupServiceSet,
        
        handlers.AuthHandlerSet,
        
        server.ProvideGraphQLHandler,
        
        routes.RouteContainerSet,

        server.New,
    )
    return &server.FiberServer{}, nil
}
