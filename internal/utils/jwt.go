package utils

import (
	"errors"
	"ganjineh-auth/internal/config"
	ent "ganjineh-auth/internal/models/entities"
	"time"

	"ganjineh-auth/pkg/ierror"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/wire"
)

type JwtPkgInterface interface {
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

type JwtPkgStruct struct {
    Conf *config.StructConfig
}

func NewJWTPkgService(cfg *config.StructConfig) *JwtPkgStruct {
    return &JwtPkgStruct{
        Conf: cfg,
    }
}

var JwtPkgSet = wire.NewSet(
    NewJWTPkgService,
    wire.Bind(new(JwtPkgInterface), new(*JwtPkgStruct)),
)

var _ JwtPkgInterface = (*JwtPkgStruct)(nil)

func (s *JwtPkgStruct) getSecret(tokenType ent.TokenType) string {
    switch tokenType {
    case ent.TokenTypeAccess:
        return s.Conf.AccessSecret
    case ent.TokenTypeRefresh:
        return s.Conf.RefreshSecret
    case ent.TokenTypeTemp:
        return s.Conf.TempSecret
    default:
        return s.Conf.AccessSecret
    }
}

func (s *JwtPkgStruct) getExpiration(tokenType ent.TokenType) time.Duration {
    switch tokenType {
    case ent.TokenTypeAccess:
        return s.Conf.AccessExpiration
    case ent.TokenTypeRefresh:
        return s.Conf.RefreshExpiration
    case ent.TokenTypeTemp:
        return s.Conf.TempExpiration
    default:
        return s.Conf.AccessExpiration
    }
}

func (s *JwtPkgStruct) generateToken(data *TokenOptions, tokenType ent.TokenType) (string, int64, *ierror.AppError) {
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
        return "", 0, ierror.NewAppError(500,1010,err.Error())
    }
    
    return tokenString, expiresAt.Unix(), nil
}

func (s *JwtPkgStruct) GenerateAccessToken(data *TokenOptions) (string, int64, *ierror.AppError) {
    return s.generateToken(data, ent.TokenTypeAccess)
}

func (s *JwtPkgStruct) GenerateRefreshToken(data *TokenOptions) (string, int64, *ierror.AppError) {
    return s.generateToken(data, ent.TokenTypeRefresh)
}

func (s *JwtPkgStruct) GenerateTempToken(data *TokenOptions) (string, int64, *ierror.AppError) {
    return s.generateToken(data, ent.TokenTypeTemp)
}

func (s *JwtPkgStruct) GenerateTokenPair(data *TokenOptions) (*ent.TokenPair, *ierror.AppError) {
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

func (s *JwtPkgStruct) ValidateAccessToken(tokenString string) (*ent.AccessTokenClaims, *ierror.AppError) {
    secret := s.getSecret(ent.TokenTypeAccess)

    token, err := jwt.ParseWithClaims(tokenString, &ent.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("error in parse jwt")
        }
        return []byte(secret), nil
    })
    if err != nil {
        return nil, ierror.NewAppError(500,1011,err.Error())
    }

    claims, ok := token.Claims.(*ent.AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, ierror.NewAppError(403,1101, "invalid access token")
	}

    if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return nil, ierror.NewAppError(403,1102, "token expired")
	}

    return claims, nil
}

func (s *JwtPkgStruct) ValidateRefreshToken(tokenString string) (*ent.RefreshTokenClaims, *ierror.AppError) {
    secret := s.getSecret(ent.TokenTypeRefresh)

    token, err := jwt.ParseWithClaims(tokenString, &ent.RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("error in parse jwt")
        }
        return []byte(secret), nil
    })
    if err != nil {
        return nil, ierror.NewAppError(500,1011,err.Error())
    }

    if claims, ok := token.Claims.(*ent.RefreshTokenClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, ierror.NewAppError(500,1101,"invalid access token")
}

func (s *JwtPkgStruct) ValidateTempToken(tokenString string) (*ent.TempTokenClamis, *ierror.AppError) {
    secret := s.getSecret(ent.TokenTypeAccess)

    token, err := jwt.ParseWithClaims(tokenString, &ent.TempTokenClamis{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("error in parse jwt")
        }
        return []byte(secret), nil
    })
    if err != nil {
        return nil, ierror.NewAppError(500,1011,err.Error())
    }

    if claims, ok := token.Claims.(*ent.TempTokenClamis); ok && token.Valid {
        return claims, nil
    }
    return nil, ierror.NewAppError(403,1101,"invalid access token")
}

