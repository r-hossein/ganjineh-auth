// middleware/provider.go
package middleware

import (
	"github.com/google/wire"
)

type MiddlewareDependencies struct {
	JWTMiddleware        JWTMiddlewareInterface
	BlackListMiddleware  BlackListMiddlewareInterface  
	PermissionMiddleware PermissionMiddlewareInterface
	RegisterMiddleware	 RegisterMiddlewareInterface
}

func NewMiddlewareDependencies(
	jwt JWTMiddlewareInterface,
	blacklist BlackListMiddlewareInterface,
	permission PermissionMiddlewareInterface,
	register	RegisterMiddlewareInterface,
) *MiddlewareDependencies {
	return &MiddlewareDependencies{
		JWTMiddleware:        jwt,
		BlackListMiddleware:  blacklist,
		PermissionMiddleware: permission,
		RegisterMiddleware: register,
	}
}

var MiddlewareSet = wire.NewSet(
	NewMiddlewareDependencies,
	MiddlewareJwtSet,
	MiddlewareBlackListSet,
	MiddlewarePermissionSet,
	MiddlewareRegisterSet,
)