package middleware

import (
	models "ganjineh-auth/internal/models/entities"
	"ganjineh-auth/internal/repositories"
	"ganjineh-auth/pkg/ierror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

type BlackListMiddleware struct {
	blackListRepo repositories.RedisBlackListRepositoryInterface
}

func NewBlackListMiddleware(blr repositories.RedisBlackListRepositoryInterface) *BlackListMiddleware {
	return &BlackListMiddleware{
		blackListRepo: blr,
	}
}

type BlackListMiddlewareInterface interface {
	Handler() fiber.Handler
}

var _ BlackListMiddlewareInterface = (*BlackListMiddleware)(nil)

var MiddlewareBlackListSet = wire.NewSet(
	NewBlackListMiddleware,
	wire.Bind(new(BlackListMiddlewareInterface), new(*BlackListMiddleware)),
)
func (m *BlackListMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sid, ok := c.Locals("sid").(string)
		if !ok || sid == "" {
					return ierror.ErrInternal
		}
		
		data, err := m.blackListRepo.GetSesion(c.Context(), sid)
		if err != nil {
			return ierror.ErrInternal
		}
		
		if data == string(models.SessionTypeRevoke) {
			return ierror.ErrTokenRevoked
		}else if data == string(models.SessionTypeUpdate){
			return ierror.ErrTokenUpdated
		}
		// Continue
		return c.Next()
	}
}
