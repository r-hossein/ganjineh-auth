package server

import (
	"fmt"
	"ganjineh-auth/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes(routeContainer *routes.RouteContainer) {
	// Apply CORS middleware
	
	fmt.Println("🚀 Registering routes...") 

	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))
	
	api := s.App.Group("/api")
	api.Get("/hello",s.HelloWorldHandler)

	routeContainer.SetupV1Routes(api)  
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}