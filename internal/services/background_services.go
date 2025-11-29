package services

import (
	"context"
	"fmt"
	"time"

	"ganjineh-auth/internal/repositories"
	"ganjineh-auth/internal/repositories/db"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgtype"
)

type BackgroundUpdateSessionParams struct {
	ID pgtype.UUID
	RefreshToken string
	ExpiresAt int64
	LastActive time.Time
}

type BackgroundServiceInterface interface {
    UpdateSession(ctx context.Context, data *BackgroundUpdateSessionParams)
	InsertError(ctx context.Context, data db.InsertErrorParams)
}

type BackgroundServiceStruct struct {
    otpRepo		repositories.RedisOTPRepositoryInterface
    userRepo    *db.Queries
}

func NewBackgroundService(
    userRepo *db.Queries,
    otpRepo repositories.RedisOTPRepositoryInterface,
) BackgroundServiceInterface {
    return &BackgroundServiceStruct{
        userRepo:   userRepo,
        otpRepo: otpRepo,
    }
}

var BackgroundServiceSet = wire.NewSet(
    NewBackgroundService,
    // wire.Bind(new(BackgroundServiceInterface), new(*AuthServiceStruct)),
)
var _ BackgroundServiceInterface = (*BackgroundServiceStruct)(nil)

func (s *BackgroundServiceStruct) UpdateSession(ctx context.Context, data *BackgroundUpdateSessionParams) {
	
	go func() {
		err := s.userRepo.UpdateSession(ctx, db.UpdateSessionParams{
		ID: data.ID,
		RefreshTokenHash: data.RefreshToken,
		RefreshTokenCreatedAt: time.Now(),
		RefreshTokenExpiresAt: time.Unix(data.ExpiresAt,0),
		LastActive: data.LastActive,
		})

		if err != nil {
			stack := "UpdateSession"
			s.InsertError(ctx,db.InsertErrorParams{
				HttpCode: 500,
				StatusCode: 1100,
				Message: err.Error(),
				StackTrace: &stack,
				Endpoint: nil,
			})
		}
	}()
}

func (s *BackgroundServiceStruct) InsertError(ctx context.Context, data db.InsertErrorParams){
	go func() {
		err := s.userRepo.InsertError(ctx,data)
		if err != nil {
			fmt.Printf("error in store error:%s",err.Error())
		}	
	}()
}
