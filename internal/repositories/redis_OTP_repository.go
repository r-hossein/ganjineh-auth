package repositories

import (
	"context"
	"encoding/json"
	"ganjineh-auth/internal/database"
	ent	"ganjineh-auth/internal/models/entities"
	"ganjineh-auth/pkg/ierror"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	StoreOTP(ctx context.Context, data *ent.OTP, expiration time.Duration) *ierror.AppError
	GetOTP(ctx context.Context, phoneNumber string) (*ent.OTP, *ierror.AppError)
	DeleteOTP(ctx context.Context, phoneNumber string) *ierror.AppError
}

type redisRepository struct {
	client *redis.Client
	prefix string
}

func NewRedisRepository(rdb database.ServiceR, prefix string) RedisRepository {
	return &redisRepository{
		client: rdb.Client(),
		prefix: prefix,
	}
}

func (r *redisRepository) getKey(phoneNumber string) string {
	return r.prefix + phoneNumber
}

func (r *redisRepository) StoreOTP(ctx context.Context, data *ent.OTP, expiration time.Duration) *ierror.AppError {

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ierror.NewAppError(1101,"can`t convert data to json!")
	}

	key := r.getKey(data.PhoneNumber)
	err = r.client.Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		return ierror.NewAppError(1301,"can`t store data in redis!")
	}

	return nil
}

func (r *redisRepository) GetOTP(ctx context.Context, phoneNumber string) (*ent.OTP, *ierror.AppError) {
	key := r.getKey(phoneNumber)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ierror.NewAppError(1102,"otp code not found!") // OTP not found
		}
		return nil, ierror.NewAppError(1301,"can`t get data from redis!")
	}

	var otpData ent.OTP
	err = json.Unmarshal([]byte(data), &otpData)
	if err != nil {
		return nil, ierror.NewAppError(1101, "can`t convert data to json!")
	}

	return &otpData, nil
}

func (r *redisRepository) DeleteOTP(ctx context.Context, phoneNumber string) *ierror.AppError {
	key := r.getKey(phoneNumber)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return ierror.NewAppError(1101,"fiald to delete data from redis!")
	}
	return nil
}
