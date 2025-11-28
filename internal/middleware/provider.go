// middleware/provider.go
package middleware

import (
	"github.com/google/wire"
)

type MiddlewareDependencies struct {
	JWTMiddleware        JWTMiddlewareInterface
	BlackListMiddleware  BlackListMiddlewareInterface  
	PermissionMiddleware PermissionMiddlewareInterface
}

func NewMiddlewareDependencies(
	jwt JWTMiddlewareInterface,
	blacklist BlackListMiddlewareInterface,
	permission PermissionMiddlewareInterface,
) *MiddlewareDependencies {
	return &MiddlewareDependencies{
		JWTMiddleware:        jwt,
		BlackListMiddleware:  blacklist,
		PermissionMiddleware: permission,
	}
}

var MiddlewareSet = wire.NewSet(
	NewMiddlewareDependencies,
	MiddlewareJwtSet,
	MiddlewareBlackListSet,
	MiddlewarePermissionSet,
)