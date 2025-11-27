package repositories

import (
	"context"
	"fmt"
	"ganjineh-auth/internal/database"
	ent "ganjineh-auth/internal/models/entities"
	"ganjineh-auth/pkg/ierror"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

type RedisBlackListRepositoryInterface interface {
	StoreBlackList(ctx context.Context, data *ent.BlackListSession) *ierror.AppError
	GetSesion(ctx context.Context, sID string) (string, *ierror.AppError)
}

type RedisBlackListRepositoryStruct struct {
	client *redis.Client
	prefix string
}

func NewRedisBlackListRepository(rdb *database.ServiceRedisStruct) RedisBlackListRepositoryInterface {
	return &RedisBlackListRepositoryStruct{
		client: rdb.Client(),
		prefix: "bl:",
	}
}

var RedisBlackListRepositorySet = wire.NewSet(
    NewRedisBlackListRepository,
    // wire.Bind(new(RedisOTPRepositoryInterface), new(*RedisOTPRepositoryStruct)),
)
var _ RedisBlackListRepositoryInterface =(*RedisBlackListRepositoryStruct)(nil)

func (r *RedisBlackListRepositoryStruct) getKey(sessionID string) string {
	return r.prefix + "session:" + sessionID
}

func (r *RedisBlackListRepositoryStruct) StoreBlackList(ctx context.Context, data *ent.BlackListSession) *ierror.AppError {
	
    key := r.getKey(data.SessionId)

    ttl := time.Until(data.ExpiresAt)
    if ttl <= 0 {
        ttl = 1 * time.Second // fallback to avoid negative TTL
    }
    
    err := r.client.Set(ctx, key, data.Status, ttl).Err()
    if err != nil {
        fmt.Printf("Redis SET error: %v\n", err.Error())
        return ierror.NewAppError(500,1301, "can't store data in redis!")
    }

    // Verify the data was stored
    val, err := r.client.Get(ctx, key).Result()
    if err != nil {
        fmt.Printf("Redis GET verification error: %v\n", err.Error())
    } else {
        fmt.Printf("Data verified in Redis: %s\n", val)
    }

    return nil
}

func (r *RedisBlackListRepositoryStruct) GetSesion(ctx context.Context, sID string) (string, *ierror.AppError) {
	key := r.getKey(sID)
	status, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ierror.NewAppError(404,1102,"otp code not found!") // OTP not found
		}
		return "", ierror.NewAppError(500,1301,"can`t get data from redis!")
	}

	return status, nil
}