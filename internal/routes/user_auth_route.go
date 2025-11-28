package routes

import (
	"ganjineh-auth/internal/handlers"
	"ganjineh-auth/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

type AuthRoutesStruct struct {
	handler handlers.AuthHandlerInterface
	middleware     *middleware.MiddlewareDependencies
}

type RouteRegistrar interface {
    RegisterRoutes(router fiber.Router)
}

func NewAuthRoutes(
	handler handlers.AuthHandlerInterface,
	middlewareDeps *middleware.MiddlewareDependencies,
) *AuthRoutesStruct {
	return &AuthRoutesStruct{
		handler: handler,
		middleware:     middlewareDeps,
	}
}

var _ RouteRegistrar = (*AuthRoutesStruct)(nil)

var AuthRoutesSet = wire.NewSet(
    NewAuthRoutes,
)
func (h *AuthRoutesStruct) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	
	otp := auth.Group("/otp")
	otp.Post("/request",h.handler.RequestOTPHandler)
	otp.Post("/verify",h.handler.VerifyOTPHandler)

	auth.Post("/register",
	h.middleware.RegisterMiddleware.Handler(),
	h.handler.RegisterUserHandler,
	)
	
	// auth.Get("/refresh",)
	
	// auth.Post("/logout",)
	// auth.Post("/logoutall",)
}
