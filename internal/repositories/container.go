// repositories/container.go
package repositories

import (
	"ganjineh-auth/internal/repositories/db"

	"github.com/google/wire"
)

type Container struct {
    OTP  RedisOTPRepositoryInterface
    User *db.Queries
    Perm RedisPermissionRepositoryInterface
    Black RedisBlackListRepositoryInterface
    // Add other repos as needed
}

func NewContainer(
    otpRepo RedisOTPRepositoryInterface,
    userRepo *db.Queries,
    blRepo RedisBlackListRepositoryInterface,
    perRepo RedisPermissionRepositoryInterface,
) *Container {
    return &Container{
        OTP:  otpRepo,
        User: userRepo,
        Perm: perRepo,
        Black: blRepo,
    }
}

var ContainerSet = wire.NewSet(
    NewContainer,
    RedisRepositorySet,  // Includes OTP repo
    QueriesSet,          // Includes user repo
    RedisBlackListRepositorySet,
    RedisPermissionRepositorySet,
)