package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	ent "ganjineh-auth/internal/models/entities"
	req "ganjineh-auth/internal/models/requests"
	res "ganjineh-auth/internal/models/responses"
	"ganjineh-auth/internal/repositories"
	"ganjineh-auth/internal/repositories/db"
	"ganjineh-auth/internal/services/otp"
	"ganjineh-auth/internal/utils"
	"ganjineh-auth/pkg/ierror"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthService interface {
    RequestOTP(ctx context.Context, phoneNumber string) (*res.OTPLoginResponse, *ierror.AppError)
    VerifyOTP(ctx context.Context, data *req.OTPVerifyRequest) (*res.OTPVerifyResponse, *ierror.AppError)
    // Register(tempToken string, userData *UserRegistration) (*AuthResult, ierror)
    // RefreshToken(refreshToken string) (*TokenPair, ierror)
    // Logout(sessionID uint, userID uuid.UUID) ierror
    // GetUserOrganizations(userID uuid.UUID) ([]Organization, ierror)
}

type authService struct {
    otpRepo		repositories.RedisRepository
	otpService	otp.OTPService
    userRepo    *db.Queries
    jwtUtil     utils.JWTService
}

func (s *authService) RequestOTP(ctx context.Context, phoneNumber string) (*res.OTPLoginResponse, *ierror.AppError) {
    
	result,erorr := s.otpService.OTPRequest(ctx, phoneNumber)

	if erorr != nil {
		return nil,erorr
	}

    return result,nil
}

func (s *authService) VerifyOTP(ctx context.Context, data *req.OTPVerifyRequest) (*res.OTPVerifyResponse, *ierror.AppError) {
    // Validate OTP
    valid, err := s.otpService.ValidateOTP(ctx,data)
    if !valid {
        return nil,err
    }
    // Check if user exists
    user, e := s.userRepo.GetUserByPhone(ctx,data.PhoneNumber)
    if e != nil {
        //if user not exist send a temp token for signup
        if errors.Is(e,pgx.ErrNoRows) {
            signToken,exp,err := s.jwtUtil.GenerateTempToken(&utils.TokenOptions{
                PhoneNumber:data.PhoneNumber,
            })
            if err != nil {
                return nil, err
            }
            return &res.OTPVerifyResponse{
                AccessToken: signToken,
                ExpiresAt: exp,
                UserExists: false,
                PhoneNumber: data.PhoneNumber,
            },nil
        }
    }
    
    if user.Status == db.UserStatusSuspended {
        return nil,ierror.NewAppError(1103,"you have baned")
    }

    if user.Status == db.UserStatusInactive{
        return nil,ierror.NewAppError(1104,"you delete your account")
    }

    // User exists - generate tokens
    sid := uuid.New()
    var result []ent.CompanyRole
    if len(user.CompanyRoles) > 0 && string(user.CompanyRoles) != "[]" {
        errr := json.Unmarshal(user.CompanyRoles, &result)
        if errr != nil {
            return nil, ierror.NewAppError(1020,errr.Error())
        }
    }

    tokens, err := s.jwtUtil.GenerateTokenPair(&utils.TokenOptions{
        UserID: user.ID.String(),
        Sid: sid.String(),
        PhoneNumber: user.PhoneNumber.(string),
        Role: user.MainRoleName.(string),
        Orgs: result,
    })
    if err != nil {
        return nil, err
    }

    ct := context.Background()

    s.userRepo.InsertUserSession(ct,db.InsertUserSessionParams{
        ID: pgtype.UUID{
            Bytes: sid, 
            Valid: true,
        },
        UserID: user.ID,
        RefreshTokenHash: tokens.RefreshToken,
        RefreshTokenCreatedAt: time.Now(),
        RefreshTokenExpiresAt: time.Unix(tokens.ExpiresIn,0),
    })

    return &res.OTPVerifyResponse{
        AccessToken: tokens.AccessToken,
        RefreshToken: tokens.RefreshToken,
        ExpiresAt: tokens.ExpiresIn,
        UserExists: true,
        FirstName: user.FirstName.(string),
        LastName: user.LastName.(string),
        UserID: user.ID.String(),
        Role: user.MainRoleName.(string),
        PhoneNumber: user.PhoneNumber.(string),
    },nil
}
