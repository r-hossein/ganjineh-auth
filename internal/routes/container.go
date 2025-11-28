// routes/container.go
package routes

import (
	"ganjineh-auth/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

type RouteContainer struct {
    registrars []RouteRegistrar
    middlewares *middleware.MiddlewareDependencies
}

func NewRouteContainer(
	middlewareDeps *middleware.MiddlewareDependencies,
	registrars ...RouteRegistrar,
) *RouteContainer {
    return &RouteContainer{
        registrars: registrars,
        middlewares: middlewareDeps,
    }
}

func (c *RouteContainer) SetupV1Routes(api fiber.Router) {
    v1 := api.Group("/v1")
    
    for _, registrar := range c.registrars {
        registrar.RegisterRoutes(v1)
    }
}

func ProvideRouteRegistrars(
    authRoutes *AuthRoutesStruct,
    graphRoutes *GraphQLRoutesStruct,
    // Add other route structs as needed
) []RouteRegistrar {
    return []RouteRegistrar{
        authRoutes,
        graphRoutes,
        // Add other route registrars here
    }
}

var RouteContainerSet = wire.NewSet(
    NewRouteContainer,
    ProvideRouteRegistrars,
    // Include all route sets
    AuthRoutesSet,
    GraphQLRoutesSet,
    // UserRoutesSet, // Add more as needed
    // AdminRoutesSet,
)