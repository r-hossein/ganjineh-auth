package repositories

import (
	"ganjineh-auth/internal/database"
	"ganjineh-auth/internal/repositories/db"
)

type Container struct {
	OTP RedisRepository
	SQLC *db.Queries
}

func NewContainer(pdb database.ServiceP, rdb database.ServiceR) *Container {
	sqlc := pdb.GetQueries()
	return &Container{
		OTP: NewRedisRepository(rdb, "otp:"),
		SQLC: sqlc,
	}
}
