package routes

import (
    "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	
    "github.com/99designs/gqlgen/graphql/handler"
    "github.com/99designs/gqlgen/graphql/playground"
    "github.com/google/wire"
)

type GraphQLRoutesStruct struct {
    graphQLHandler *handler.Server
}

func NewGraphQLRoutes(graphQLHandler *handler.Server) *GraphQLRoutesStruct {
    return &GraphQLRoutesStruct{
        graphQLHandler: graphQLHandler,
    }
}

var _ RouteRegistrar = (*GraphQLRoutesStruct)(nil)

var GraphQLRoutesSet = wire.NewSet(
    NewGraphQLRoutes,
)

func (g *GraphQLRoutesStruct) RegisterRoutes(router fiber.Router) {
    // GraphQL endpoint
    router.All("/graph", adaptor.HTTPHandler(g.graphQLHandler) )
    
    // GraphQL Playground
    router.Get("/playground", adaptor.HTTPHandler(
        playground.Handler("GraphQL Playground", "/api/v1/graph"),
    ))
}