package repositories

import (
	"ganjineh-auth/internal/repositories/db"
	"ganjineh-auth/internal/database"
	"github.com/google/wire"
)

// type Container struct {
// 	OTP RedisRepository
// 	SQLC *db.Queries
// }

func NewQueries(dbp *database.ServicePostgresStruct) *db.Queries { // Add * here
    return db.New(dbp.Pool)
}

var QueriesSet = wire.NewSet(NewQueries)