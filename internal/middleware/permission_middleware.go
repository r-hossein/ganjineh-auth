package middleware

import (
	models "ganjineh-auth/internal/models/entities"
	"ganjineh-auth/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

type PermissionMiddleware struct {
	perRepo repositories.RedisPermissionRepositoryInterface
}

func NewPermissionMiddleware(per repositories.RedisPermissionRepositoryInterface) *PermissionMiddleware {
	return &PermissionMiddleware{
		perRepo: per,
	}
}

type PermissionMiddlewareInterface interface {
	Handler() fiber.Handler
}

var _ PermissionMiddlewareInterface = (*PermissionMiddleware)(nil)

var MiddlewarePermissionSet = wire.NewSet(
	NewPermissionMiddleware,
	wire.Bind(new(PermissionMiddlewareInterface), new(*PermissionMiddleware)),
)

func (m *PermissionMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role_main").(string)
		data, err := m.perRepo.GetPermission(c.Context(), role)
		if err != nil {
			return err
		}
		
		c.Locals("permission_main",data.PermissionCodes)
		
		orgsAny := c.Locals("organizations")
		if orgsAny == nil{
			return c.Next()
		}
		
		orgs, ok := orgsAny.([]models.CompanyRole)
        if !ok || len(orgs) == 0 {
            // Organizations exist but empty or wrong type → continue
            return c.Next()
        }
		
        for i:= range orgs{
        	roleName := orgs[i].RoleName
        	perm, err := m.perRepo.GetPermission(c.Context(), roleName)
			if err != nil {
				return err
			}
			orgs[i].Permissions=perm.PermissionCodes
        }
		c.Locals("organizations", orgs)
		// Continue
		return c.Next()
	}
}
