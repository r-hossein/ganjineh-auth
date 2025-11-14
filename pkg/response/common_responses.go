package responses

import (
    "ganjineh-auth/pkg/ierror"
)

type BaseResponse struct {
    Success bool        `json:"success"`
    Code    int         `json:"code"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
    Message string `json:"message"`
}

func SuccessResponse(data interface{}, code int) BaseResponse {
    return BaseResponse{
        Success: true,
        Code: code,
        Data:    data,
    }
}

func SuccessWithMessage(message string, data interface{}, code int) BaseResponse {
    return BaseResponse{
        Success: true,
        Code: code,
        Message: message,
        Data:    data,
    }
}

func ErrorResponse(errors *ierror.AppError) BaseResponse {
    return BaseResponse{
        Success: false,
        Code: errors.Code,
        Error: &ErrorInfo{
            Message: errors.Message,
        },
    }
}