package middleware

import (

	"ganjineh-auth/internal/utils"
	"ganjineh-auth/pkg/ierror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

type RefreshtokenMiddleware struct {
	JWT utils.JwtPkgInterface
}

func NewRefreshtokenMiddleware(jwt utils.JwtPkgInterface) *RefreshtokenMiddleware {
	return &RefreshtokenMiddleware{
		JWT: jwt,
	}
}

type RefreshtokenMiddlewareInterface interface {
	Handler() fiber.Handler
}

var _ RefreshtokenMiddlewareInterface = (*RefreshtokenMiddleware)(nil)

var MiddlewareRefreshtokenSet = wire.NewSet(
	NewRefreshtokenMiddleware,
	wire.Bind(new(RefreshtokenMiddlewareInterface), new(*RefreshtokenMiddleware)),
)

func (m *RefreshtokenMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		refreshToken := c.Cookies("refresh_token")
		if refreshToken == "" {
			return ierror.ErrTokenInvalid
		}
		// Validate temp Token using your JWT service
		claims, err := m.JWT.ValidateRefreshToken(refreshToken)
		if err != nil {
			// err is *ierror.AppError
			return err
		}
		// Attach claims to context
		c.Locals("token",refreshToken)
		c.Locals("sid", claims.Sid)
		c.Locals("phone_number",claims.PhoneNumber)
		c.Locals("role_main", claims.Role)
		c.Locals("organizations", claims.Organizations)
		c.Locals("userID", claims.Subject)

		// Continue
		return c.Next()
	}
}
