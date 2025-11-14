package ierror

import (
    "net/http"
)

type AppError struct {
    Code    int
    Message string
}

func (e *AppError) Error() string {
    return e.Message
}

func NewAppError(code int, message string) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
    }
}

// خطاهای رایج
var (
    ErrNotFound     = NewAppError(http.StatusNotFound, "not found!")
    ErrUnauthorized = NewAppError(http.StatusUnauthorized, "دسترسی غیرمجاز")
    ErrBadRequest   = NewAppError(http.StatusBadRequest, "invalid request")
    ErrInternal     = NewAppError(http.StatusInternalServerError, "خطای داخلی سرور")
    
    // خطاهای دامنه (Domain Errors)
    ErrUserNotFound      = NewAppError(http.StatusNotFound, "کاربر یافت نشد")
    ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, "ایمیل یا رمز عبور اشتباه است")
    ErrEmailExists       = NewAppError(http.StatusConflict, "این ایمیل قبلاً ثبت شده است")
    ErrInvalidOTP        = NewAppError(http.StatusBadRequest, "کد OTP نامعتبر است")
    ErrOTPExpired        = NewAppError(http.StatusBadRequest, "کد OTP منقضی شده است")
)