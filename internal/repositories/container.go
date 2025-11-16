// repositories/container.go
package repositories

import (
	"ganjineh-auth/internal/repositories/db"

	"github.com/google/wire"
)

type Container struct {
    OTP  RedisOTPRepositoryInterface
    User *db.Queries
    // Add other repos as needed
}

func NewContainer(
    otpRepo RedisOTPRepositoryInterface,
    userRepo *db.Queries,
) *Container {
    return &Container{
        OTP:  otpRepo,
        User: userRepo,
    }
}

var ContainerSet = wire.NewSet(
    NewContainer,
    RedisRepositorySet,  // Includes OTP repo
    QueriesSet,          // Includes user repo
)