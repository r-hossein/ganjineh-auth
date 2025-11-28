package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"ganjineh-auth/internal/database"
	ent "ganjineh-auth/internal/models/entities"
	"ganjineh-auth/pkg/ierror"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

type RedisOTPRepositoryInterface interface {
	StoreOTP(ctx context.Context, data *ent.OTP, expiration time.Duration) *ierror.AppError
	GetOTP(ctx context.Context, phoneNumber string) (*ent.OTP, *ierror.AppError)
	DeleteOTP(ctx context.Context, phoneNumber string) *ierror.AppError
}

type RedisOTPRepositoryStruct struct {
	client *redis.Client
	prefix string
}

func NewRedisRepository(rdb *database.ServiceRedisStruct) RedisOTPRepositoryInterface {
	return &RedisOTPRepositoryStruct{
		client: rdb.Client(),
		prefix: "otp:",
	}
}

var RedisRepositorySet = wire.NewSet(
    NewRedisRepository,
    // wire.Bind(new(RedisOTPRepositoryInterface), new(*RedisOTPRepositoryStruct)),
)
var _ RedisOTPRepositoryInterface =(*RedisOTPRepositoryStruct)(nil)

func (r *RedisOTPRepositoryStruct) getKey(phoneNumber string) string {
	return r.prefix + phoneNumber
}

func (r *RedisOTPRepositoryStruct) StoreOTP(ctx context.Context, data *ent.OTP, expiration time.Duration) *ierror.AppError {

    jsonData, err := json.Marshal(data)
    if err != nil {
        fmt.Printf("JSON Marshal error: %v\n", err.Error())
        return ierror.NewAppError(500,1101, "can't convert data to json!")
    }
    
    fmt.Printf("JSON data: %s\n", string(jsonData))

    key := r.getKey(data.PhoneNumber)
    fmt.Printf("Redis key: %s\n", key)

    err = r.client.Set(ctx, key, jsonData, expiration).Err()
    if err != nil {
        fmt.Printf("Redis SET error: %v\n", err.Error())
        return ierror.NewAppError(500,1301, "can't store data in redis!")
    }

    // Verify the data was stored
    // val, err := r.client.Get(ctx, key).Result()
    // if err != nil {
    //     return  ierror.NewAppError(500,1301,err.Error())
    // } else {
    //     return nil
    // }
	return  nil
}

func (r *RedisOTPRepositoryStruct) GetOTP(ctx context.Context, phoneNumber string) (*ent.OTP, *ierror.AppError) {
	key := r.getKey(phoneNumber)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ierror.NewAppError(404,1102,"otp code not found!") // OTP not found
		}
		return nil, ierror.NewAppError(500,1301,"can`t get data from redis!")
	}

	var otpData ent.OTP
	err = json.Unmarshal([]byte(data), &otpData)
	if err != nil {
		return nil, ierror.NewAppError(500,1101, "can`t convert data to json!")
	}

	return &otpData, nil
}

func (r *RedisOTPRepositoryStruct) DeleteOTP(ctx context.Context, phoneNumber string) *ierror.AppError {
	key := r.getKey(phoneNumber)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return ierror.NewAppError(500,1101,"fiald to delete data from redis!")
	}
	return nil
}
