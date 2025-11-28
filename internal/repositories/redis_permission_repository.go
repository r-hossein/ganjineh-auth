package repositories

import (
	"context"
	"fmt"
	"encoding/json"
	"ganjineh-auth/internal/database"
	ent "ganjineh-auth/internal/models/entities"
	"ganjineh-auth/pkg/ierror"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

type RedisPermissionRepositoryInterface interface {
	StorePermission(ctx context.Context, data *ent.Role) *ierror.AppError
	GetPermission(ctx context.Context, roleName string) (*ent.Role, *ierror.AppError)
	ClearAllRoles(ctx context.Context) *ierror.AppError 
}

type RedisPermissionRepositoryStruct struct {
	client *redis.Client
	prefix string
}

func NewRedisPermissionRepository(rdb *database.ServiceRedisStruct) RedisPermissionRepositoryInterface {
	return &RedisPermissionRepositoryStruct{
		client: rdb.Client(),
		prefix: "role:",
	}
}

var RedisPermissionRepositorySet = wire.NewSet(
    NewRedisPermissionRepository,
    // wire.Bind(new(RedisOTPRepositoryInterface), new(*RedisOTPRepositoryStruct)),
)
var _ RedisPermissionRepositoryInterface =(*RedisPermissionRepositoryStruct)(nil)

func (r *RedisPermissionRepositoryStruct) getKey(RoleName string) string {
	return r.prefix + RoleName +":permissions"
}

func (r *RedisPermissionRepositoryStruct) StorePermission(ctx context.Context, data *ent.Role) *ierror.AppError {
 	jsonData, err := json.Marshal(data)
    if err != nil {
        fmt.Printf("JSON Marshal error: %v\n", err.Error())
        return ierror.NewAppError(500,1101, "can't convert data to json!")
    }
	
    key := r.getKey(data.Name)
    
    err = r.client.SAdd(ctx, key, jsonData).Err()
    if err != nil {
        fmt.Printf("Redis SET error: %v\n", err.Error())
        return ierror.NewAppError(500,1301, "can't store data in redis!")
    }

    return nil
}

func (r *RedisPermissionRepositoryStruct) GetPermission(ctx context.Context, roleName string) (*ent.Role, *ierror.AppError) {
	key := r.getKey(roleName)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ierror.NewAppError(404,1102,"otp code not found!") // OTP not found
		}
		return nil, ierror.NewAppError(500,1301,"can`t get data from redis!")
	}

	var roleData ent.Role
	err = json.Unmarshal([]byte(data), &roleData)
	if err != nil {
		return nil, ierror.NewAppError(500,1101, "can`t convert data to json!")
	}
	
	return &roleData, nil
}

func (r *RedisPermissionRepositoryStruct) ClearAllRoles(ctx context.Context) *ierror.AppError {
    key := r.prefix+"*:permissions"
	iter := r.client.Scan(ctx, 0, key, 0).Iterator()

    for iter.Next(ctx) {
        if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
            return ierror.NewAppError(500,1304,"error in rewite role")
        }
    }

    return nil
}