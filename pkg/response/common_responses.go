package responses

import (
    "ganjineh-auth/pkg/ierror"
)

type BaseResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
    Message string `json:"message"`
    Reason  string `json:"reason,omitempty"`
}

func SuccessResponse(data interface{}, code int) BaseResponse {
    return BaseResponse{
        Code: code,
        Data:    data,
    }
}

func SuccessWithMessage(message string, data interface{}, code int) BaseResponse {
    return BaseResponse{
        Code: code,
        Message: message,
        Data:    data,
    }
}

func ErrorResponse(err error) BaseResponse {

    // Handle AuthError
    if authErr, ok := err.(*ierror.AuthError); ok {
        return BaseResponse{
            Code: authErr.Code,
            Error: &ErrorInfo{
                Message: authErr.Message,
                Reason:  authErr.Reason,
            },
        }
    }

    // Handle AppError (your existing errors)
    if appErr, ok := err.(*ierror.AppError); ok {
        return BaseResponse{
            Code: appErr.Code,
            Error: &ErrorInfo{
                Message: appErr.Message,
            },
        }
    }

    // Unknown error → internal error
    return BaseResponse{
        Code: 500,
        Error: &ErrorInfo{
            Message: "internal server error",
        },
    }
}
