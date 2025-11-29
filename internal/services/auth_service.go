package services

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
	"ganjineh-auth/internal/utils"
	"ganjineh-auth/pkg/ierror"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthServiceInterface interface {
    RequestOTP(ctx context.Context, phoneNumber string) (*res.OTPLoginResponse, *ierror.AppError)
    VerifyOTP(ctx context.Context, data *req.OTPVerifyRequest) (*res.OTPVerifyResponse, *ierror.AppError)
    Register(ctx context.Context, data *req.RegisterUserFirstRequest) (*res.OTPVerifyResponse, *ierror.AppError)
    RefreshToken(ctx context.Context) (*res.RefreshToken, *ierror.AppError)
    // Logout(sessionID uint, userID uuid.UUID) ierror
    // GetUserOrganizations(userID uuid.UUID) ([]Organization, ierror)
}

type AuthServiceStruct struct {
    otpRepo		repositories.RedisOTPRepositoryInterface
    userRepo    *db.Queries
    jwtUtil     utils.JwtPkgInterface
    otpServ     OTPServiceInterface
    bgService   BackgroundServiceInterface
}

func NewAuthService(
    otpService OTPServiceInterface,
    userRepo *db.Queries,
    jwtUtil utils.JwtPkgInterface,
    otpRepo repositories.RedisOTPRepositoryInterface,
    bgServ  BackgroundServiceInterface,
) AuthServiceInterface {
    return &AuthServiceStruct{
        otpServ: otpService,
        userRepo:   userRepo,
        jwtUtil:    jwtUtil,
        otpRepo: otpRepo,
        bgService: bgServ,
    }
}

var AuthServiceSet = wire.NewSet(
    NewAuthService,
    // wire.Bind(new(AuthServiceInterface), new(*AuthServiceStruct)),
)
var _ AuthServiceInterface = (*AuthServiceStruct)(nil)

func (s *AuthServiceStruct) RequestOTP(ctx context.Context, phoneNumber string) (*res.OTPLoginResponse, *ierror.AppError) {
    
	result,erorr := s.otpServ.OTPRequest(ctx, phoneNumber)

	if erorr != nil {
		return nil,erorr
	}

    return result,nil
}

func (s *AuthServiceStruct) VerifyOTP(ctx context.Context, data *req.OTPVerifyRequest) (*res.OTPVerifyResponse, *ierror.AppError) {
    // Validate OTP
    valid, err := s.otpServ.ValidateOTP(ctx,data)
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
        return nil,ierror.NewAppError(403,1103,"you have baned")
    }

    if user.Status == db.UserStatusInactive{
        return nil,ierror.NewAppError(403,1104,"you delete your account")
    }

    // User exists - generate tokens
    sid := uuid.New()
    var result []ent.CompanyRole
    if len(user.CompanyRoles) > 0 && string(user.CompanyRoles) != "[]" {
        errr := json.Unmarshal(user.CompanyRoles, &result)
        if errr != nil {
            return nil, ierror.NewAppError(500,1020,errr.Error())
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

func (s *AuthServiceStruct) Register(ctx context.Context, data *req.RegisterUserFirstRequest) (*res.OTPVerifyResponse, *ierror.AppError){

    phoneNumber := ctx.Value("phone_number").(string)
    
    user, err := s.userRepo.CreateUser(ctx, db.CreateUserParams{
        PhoneNumber: phoneNumber,
        FirstName: data.FirstName,
        LastName: data.LastName,
        Gender: db.GenderEnum(data.Gender),
        RoleID: 3,
    })

    if err != nil {
        return nil, ierror.NewAppError(400,1100,err.Error())
    }

    role, err := s.userRepo.GetRoleByID(ctx,user.RoleID)
    if err != nil {
        return nil, ierror.NewAppError(400,1100,err.Error())
    }

    sid := uuid.New()
    tokens, er := s.jwtUtil.GenerateTokenPair(&utils.TokenOptions{
        UserID: user.ID.String(),
        Sid: sid.String(),
        PhoneNumber: user.PhoneNumber,
        Role: role,
        Orgs: []ent.CompanyRole{},
    })
    if er != nil {
        return nil, er
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
        FirstName: user.FirstName,
        LastName: user.LastName,
        UserID: user.ID.String(),
        Role: role,
        PhoneNumber: user.PhoneNumber,
    },nil

}

func (s *AuthServiceStruct) RefreshToken(ctx context.Context) (*res.RefreshToken, *ierror.AppError)  {
    
    phoneNumber, ok := ctx.Value("phone_number").(string)
    if !ok {
        return nil, ierror.NewAppError(500, 1100, "invalid phone number")
    }
    
    refreshToken, ok := ctx.Value("token").(string)
    if !ok {
        return nil, ierror.NewAppError(500, 1100, "invalid refresh token")
    }
    
    sid, ok := ctx.Value("sid").(string)
    if !ok {
        return nil, ierror.NewAppError(500, 1100, "invalid session id")
    }
    
    role, ok := ctx.Value("role_main").(string)
    if !ok {
        return nil, ierror.NewAppError(500, 1100, "invalid role")
    }
    
    orgs, ok := ctx.Value("organizations").([]ent.CompanyRole)
    if !ok {
        return nil, ierror.NewAppError(500, 1100, "invalid organizations")
    }
    
    userId, ok := ctx.Value("userID").(string)
    if !ok {
        return nil, ierror.NewAppError(500, 1100, "invalid user id")
    }

    // Convert string IDs to UUID
    sessionUUID, err := uuid.Parse(sid)
    if err != nil {
        return nil, ierror.NewAppError(500, 1100, "invalid session id format")
    }
    
    userUUID, err := uuid.Parse(userId)
    if err != nil {
        return nil, ierror.NewAppError(500, 1100, "invalid user id format")
    }
    
    session, err := s.userRepo.GetSession(ctx,db.GetSessionParams{
        ID: pgtype.UUID{
            Bytes: sessionUUID, 
            Valid: true,
        },
        UserID: pgtype.UUID{
            Bytes: userUUID, 
            Valid: true,
        },
    })

    if err != nil {
        return nil, ierror.NewAppError(500,1303,err.Error())
    }
    if session.Status == db.SessionStatusTypeRevoked {
        return nil, ierror.NewAppError(403,1104,"invalid refreshtoken")
    }
    if session.RefreshTokenHash != refreshToken {
        return nil, ierror.NewAppError(403,1102,"invalid refreshtoken")
    }else if time.Now().After(session.RefreshTokenExpiresAt) {
        return nil, ierror.NewAppError(403,1103, "token expired")
    }

    data , er := s.jwtUtil.GenerateTokenPair(&utils.TokenOptions{
        UserID: userId,
        Sid: sid,
        PhoneNumber: phoneNumber,
        Role: role,
        Orgs: orgs,
    })
    if er != nil {
        return nil, er
    }
    
    s.bgService.UpdateSession(context.Background(), &BackgroundUpdateSessionParams{
        ID: session.ID,
        RefeshToken: data.RefreshToken,
        ExpiresAt: data.ExpiresIn,
        LastActive: time.Now(),
    })
    
    return &res.RefreshToken{
        RefreshToken: data.RefreshToken,
        AccessToken: data.AccessToken,
        ExpiresAt: data.ExpiresIn,
    } , nil
}
