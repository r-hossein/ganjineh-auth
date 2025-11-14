package utils

import (
	"errors"
	"ganjineh-auth/internal/config"
	ent "ganjineh-auth/internal/models/entities"
	"time"

	"ganjineh-auth/pkg/ierror"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
    GenerateAccessToken(data *TokenOptions) (string, int64, *ierror.AppError)
    GenerateRefreshToken(data *TokenOptions) (string, int64, *ierror.AppError)
    GenerateTempToken(data *TokenOptions) (string, int64, *ierror.AppError)
    GenerateTokenPair(data *TokenOptions) (*ent.TokenPair, *ierror.AppError)
    ValidateAccessToken(tokenString string) (*ent.AccessTokenClaims, *ierror.AppError)
    ValidateRefreshToken(tokenString string) (*ent.RefreshTokenClaims, *ierror.AppError)
    ValidateTempToken(tokenString string) (*ent.TempTokenClamis, *ierror.AppError)
}

type TokenOptions struct {
    UserID      string
    Sid         string
    PhoneNumber string
    Role        string
    Orgs        []ent.CompanyRole
}

type jwtService struct {
    config *config.JWTConfig
}

func NewJWTService(cfg *config.JWTConfig) JWTService {
    return &jwtService{
        config: cfg,
    }
}

func (s *jwtService) getSecret(tokenType ent.TokenType) string {
    switch tokenType {
    case ent.TokenTypeAccess:
        return s.config.AccessSecret
    case ent.TokenTypeRefresh:
        return s.config.RefreshSecret
    case ent.TokenTypeTemp:
        return s.config.TempSecret
    default:
        return s.config.AccessSecret
    }
}

func (s *jwtService) getExpiration(tokenType ent.TokenType) time.Duration {
    switch tokenType {
    case ent.TokenTypeAccess:
        return s.config.AccessExpiration
    case ent.TokenTypeRefresh:
        return s.config.RefreshExpiration
    case ent.TokenTypeTemp:
        return s.config.TempExpiration
    default:
        return s.config.AccessExpiration
    }
}

func (s *jwtService) generateToken(data *TokenOptions, tokenType ent.TokenType) (string, int64, *ierror.AppError) {
    expiration := s.getExpiration(tokenType)
    expiresAt := time.Now().Add(expiration)
    var claims jwt.Claims
    switch tokenType {
    case ent.TokenTypeAccess:
        claims =&ent.AccessTokenClaims{
            Sid: data.Sid,
            PhoneNumber: data.PhoneNumber,
            Role: data.Role,
            TokenType: tokenType,
            Organizations: data.Orgs,
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(expiresAt),
                IssuedAt: jwt.NewNumericDate(time.Now()),
                Subject: data.UserID,
                ID: data.Sid,
            },
        }
    case ent.TokenTypeRefresh:
        claims =&ent.RefreshTokenClaims{
            Sid: data.Sid,
            PhoneNumber: data.PhoneNumber,
            Role: data.Role,
            TokenType: tokenType,
            Organizations: data.Orgs,
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(expiresAt),
                IssuedAt: jwt.NewNumericDate(time.Now()),
                Subject: data.UserID,
                ID: data.Sid,
            },
        }
    case ent.TokenTypeTemp :
        claims =&ent.TempTokenClamis{
            PhoneNumber: data.PhoneNumber,
            TokenType: tokenType,
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(expiresAt),
                IssuedAt: jwt.NewNumericDate(time.Now()),
            },
        }
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    secret := s.getSecret(tokenType)
    
    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", 0, ierror.NewAppError(1010,err.Error())
    }
    
    return tokenString, expiresAt.Unix(), nil
}

func (s *jwtService) GenerateAccessToken(data *TokenOptions) (string, int64, *ierror.AppError) {
    return s.generateToken(data, ent.TokenTypeAccess)
}

func (s *jwtService) GenerateRefreshToken(data *TokenOptions) (string, int64, *ierror.AppError) {
    return s.generateToken(data, ent.TokenTypeRefresh)
}

func (s *jwtService) GenerateTempToken(data *TokenOptions) (string, int64, *ierror.AppError) {
    return s.generateToken(data, ent.TokenTypeTemp)
}

func (s *jwtService) GenerateTokenPair(data *TokenOptions) (*ent.TokenPair, *ierror.AppError) {
    accessToken, _, err := s.GenerateAccessToken(data)
    if err != nil {
        return nil, err
    }
    
    refreshToken, refreshExp, err := s.GenerateRefreshToken(data)
    if err != nil {
        return nil, err
    }
    
    return &ent.TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    refreshExp,
    }, nil
}

func (s *jwtService) ValidateAccessToken(tokenString string) (*ent.AccessTokenClaims, *ierror.AppError) {
    secret := s.getSecret(ent.TokenTypeAccess)

    token, err := jwt.ParseWithClaims(tokenString, &ent.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("error in parse jwt")
        }
        return []byte(secret), nil
    })
    if err != nil {
        return nil, ierror.NewAppError(1011,err.Error())
    }

    if claims, ok := token.Claims.(*ent.AccessTokenClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, ierror.NewAppError(1101,"invalid access token")
}

func (s *jwtService) ValidateRefreshToken(tokenString string) (*ent.RefreshTokenClaims, *ierror.AppError) {
    secret := s.getSecret(ent.TokenTypeRefresh)

    token, err := jwt.ParseWithClaims(tokenString, &ent.RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("error in parse jwt")
        }
        return []byte(secret), nil
    })
    if err != nil {
        return nil, ierror.NewAppError(1011,err.Error())
    }

    if claims, ok := token.Claims.(*ent.RefreshTokenClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, ierror.NewAppError(1101,"invalid access token")
}

func (s *jwtService) ValidateTempToken(tokenString string) (*ent.TempTokenClamis, *ierror.AppError) {
    secret := s.getSecret(ent.TokenTypeAccess)

    token, err := jwt.ParseWithClaims(tokenString, &ent.TempTokenClamis{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("error in parse jwt")
        }
        return []byte(secret), nil
    })
    if err != nil {
        return nil, ierror.NewAppError(1011,err.Error())
    }

    if claims, ok := token.Claims.(*ent.TempTokenClamis); ok && token.Valid {
        return claims, nil
    }
    return nil, ierror.NewAppError(1101,"invalid access token")
}

