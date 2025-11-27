package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"ganjineh-auth/internal/utils"
	"github.com/google/wire"
)

type BlackListMiddleware struct {
	JWT utils.JwtPkgInterface
}

func NewBlackListMiddleware(jwt utils.JwtPkgInterface) *BlackListMiddleware {
	return &BlackListMiddleware{
		JWT: jwt,
	}
}

var MiddlewareblackListSet = wire.NewSet(
    NewBlackListMiddleware,
)

func (m *BlackListMiddleware) Handler() fiber.Handler {
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

		// Validate Access Token using your JWT service
		claims, err := m.JWT.ValidateAccessToken(tokenString)
		if err != nil {
			// err is *ierror.AppError
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		// Attach claims to context
		c.Locals("claims", claims)
		c.Locals("userID", claims.Subject)
		c.Locals("role", claims.Role)

		// Continue
		return c.Next()
	}
}
