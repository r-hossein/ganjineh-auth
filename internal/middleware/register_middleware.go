package middleware

import (
	"strings"

	"ganjineh-auth/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

type RegisterMiddleware struct {
	JWT utils.JwtPkgInterface
}

func NewRegisterMiddleware(jwt utils.JwtPkgInterface) *RegisterMiddleware {
	return &RegisterMiddleware{
		JWT: jwt,
	}
}

type RegisterMiddlewareInterface interface {
	Handler() fiber.Handler
}

var _ RegisterMiddlewareInterface = (*RegisterMiddleware)(nil)

var MiddlewareRegisterSet = wire.NewSet(
	NewRegisterMiddleware,
	wire.Bind(new(RegisterMiddlewareInterface), new(*RegisterMiddleware)),
)

func (m *RegisterMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		// Expect "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.ErrUnauthorized
		}

		tokenString := parts[1]

		// Validate temp Token using your JWT service
		claims, err := m.JWT.ValidateTempToken(tokenString)
		if err != nil {
			// err is *ierror.AppError
			return err
		}
		// Attach claims to context
		c.Locals("phone_number",claims.PhoneNumber)

		// Continue
		return c.Next()
	}
}
