package routes

import (
	"ganjineh-auth/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

type AuthRoutesStruct struct {
	handler handlers.AuthHandlerInterface
}

type RouteRegistrar interface {
    RegisterRoutes(router fiber.Router)
}

func NewAuthRoutes(handler handlers.AuthHandlerInterface) *AuthRoutesStruct {
	return &AuthRoutesStruct{
		handler: handler,
	}
}

var _ RouteRegistrar = (*AuthRoutesStruct)(nil)

var AuthRoutesSet = wire.NewSet(
    NewAuthRoutes,
    wire.Bind(new(RouteRegistrar), new(*AuthRoutesStruct)), // Bind to interface
)
func (h *AuthRoutesStruct) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	
	otp := auth.Group("/otp")
	otp.Post("/request",h.handler.RequestOTPHandler)
	otp.Post("/verify",h.handler.VerifyOTPHandler)

	// auth.Post("/register",)
	
	// auth.Get("/refresh",)
	
	// auth.Post("/logout",)
	// auth.Post("/logoutall",)
}
