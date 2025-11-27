package ierror

import (
    "net/http"
)

type AppError struct {
    HttpStatus int
    Code    int
    Message string
}

func (e *AppError) Error() string {
    return e.Message
}

func NewAppError(httpStatus int, code int, message string) *AppError {
    return &AppError{
        HttpStatus: httpStatus,
        Code:    code,
        Message: message,
    }
}

// خطاهای رایج
var (
    ErrNotFound     = NewAppError(http.StatusNotFound, 404 ,"not found!")
    ErrUnauthorized = NewAppError(http.StatusUnauthorized, 403 ,"دسترسی غیرمجاز")
    ErrBadRequest   = NewAppError(http.StatusBadRequest, 404 ,"invalid request")
    ErrInternal     = NewAppError(http.StatusInternalServerError,404, "خطای داخلی سرور")
    
    // خطاهای دامنه (Domain Errors)
    ErrUserNotFound      = NewAppError(http.StatusNotFound, 404 ,"کاربر یافت نشد")
    ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, 404,"ایمیل یا رمز عبور اشتباه است")
    ErrEmailExists       = NewAppError(http.StatusConflict,404, "این ایمیل قبلاً ثبت شده است")
    ErrInvalidOTP        = NewAppError(http.StatusBadRequest,404, "کد OTP نامعتبر است")
    ErrOTPExpired        = NewAppError(http.StatusBadRequest,404, "کد OTP منقضی شده است")
)