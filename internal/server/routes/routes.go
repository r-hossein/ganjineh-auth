package routes

import (
	"github.com/gofiber/fiber/v2"
	"ganjineh-auth/internal/handlers/v1"
)

func SetupV1Routes(api fiber.Router) {
	V1 := api.Group("/v1")

	authRoutes := NewAuthRoutes(&v1.AuthHandler{})
	authRoutes.RegisterUserAuthRoutes(V1)
	
}

