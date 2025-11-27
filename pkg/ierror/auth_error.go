package ierror

type AuthError struct {
    HttpStatus int
    Code       int
    Message    string
    Reason     string
}

func (e *AuthError) Error() string {
    return e.Message
}

func NewAuthError(httpStatus, code int, msg, reason string) *AuthError {
    return &AuthError{
        HttpStatus: httpStatus,
        Code:       code,
        Message:    msg,
        Reason:     reason,
    }
}

var (
    // Token errors
    ErrTokenInvalid = NewAuthError(401, 1101, "invalid token", "auth_required")
    ErrTokenExpired = NewAuthError(401, 1102, "access token expired", "refresh_required")
    ErrTokenRevoked = NewAuthError(401, 1201, "token revoked", "revoked")

    // User account errors
    ErrUserBanned = NewAuthError(401, 1202, "user banned", "banned")
)