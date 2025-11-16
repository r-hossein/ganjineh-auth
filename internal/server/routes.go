package server

import (
	// "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"ganjineh-auth/internal/routes"
)

func (s *FiberServer) RegisterFiberRoutes(routeContainer *routes.RouteContainer) {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	api := s.App.Group("/api")
	routeContainer.SetupV1Routes(api)  
}
