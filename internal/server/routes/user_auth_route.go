package routes

import (
	"ganjineh-auth/internal/handlers/v1"

	"github.com/gofiber/fiber/v2"
)

type AuthRoutes struct {
	Handler v1.AuthHandlerInterface
}

func NewAuthRoutes(handler v1.AuthHandlerInterface) *AuthRoutes {
	return &AuthRoutes{
		Handler: handler,
	}
}

func (h *AuthRoutes) RegisterUserAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	
	otp := auth.Group("/otp")
	otp.Post("/request",h.Handler.RequestOTPHandler)
	otp.Post("/verify",h.Handler.VerifyOTPHandler)

	// auth.Post("/register",)
	
	// auth.Get("/refresh",)
	
	// auth.Post("/logout",)
	// auth.Post("/logoutall",)
}
