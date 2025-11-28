package routes

import (
	"ganjineh-auth/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/google/wire"
)

type GraphQLRoutesStruct struct {
    graphQLHandler *handler.Server
    middleware     *middleware.MiddlewareDependencies
}

func NewGraphQLRoutes(
	graphQLHandler *handler.Server,
	middlewareDeps *middleware.MiddlewareDependencies,
) *GraphQLRoutesStruct {
    return &GraphQLRoutesStruct{
        graphQLHandler: graphQLHandler,
        middleware:     middlewareDeps,
    }
}

var _ RouteRegistrar = (*GraphQLRoutesStruct)(nil)

var GraphQLRoutesSet = wire.NewSet(
    NewGraphQLRoutes,
)

func (g *GraphQLRoutesStruct) RegisterRoutes(router fiber.Router) {
    // GraphQL endpoint
    router.All("/graphql",
    	g.middleware.JWTMiddleware.Handler(),
     	g.middleware.BlackListMiddleware.Handler(),
      	g.middleware.PermissionMiddleware.Handler(),
    	adaptor.HTTPHandler(g.graphQLHandler),
    )
    
    // GraphQL Playground
    router.Get("/playground", adaptor.HTTPHandler(
        playground.Handler("GraphQL Playground", "/api/v1/graphql"),
    ))
}