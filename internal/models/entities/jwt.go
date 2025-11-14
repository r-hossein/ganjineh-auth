package models

import "github.com/golang-jwt/jwt/v5"

type TokenType string

const (
    TokenTypeAccess  TokenType = "access"
    TokenTypeRefresh TokenType = "refresh"
    TokenTypeTemp    TokenType = "temp"
)

type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
}

type TempTokenClamis struct {
    PhoneNumber     string    `json:"phone_number"`
    TokenType   TokenType `json:"token_type"`
    jwt.RegisteredClaims
}

type AccessTokenClaims struct {
    Sid     string    `json:"sid"`
    PhoneNumber     string    `json:"phone_number"`
    Role    string    `json:"role"`
    TokenType   TokenType `json:"token_type"`
    Organizations []CompanyRole `json:"organizations"`
    jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
    Sid string  `json:"sid"`
    PhoneNumber string    `json:"phone_number"`
    Role string    `json:"role"`
    TokenType   TokenType `json:"token_type"`
    Organizations []CompanyRole `json:"organizations"`
    jwt.RegisteredClaims
}

type CompanyRole struct {
	CompanyID string `json:"company_id"`
	RoleName  string `json:"role_name"`
}