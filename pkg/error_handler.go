package pkg

import (
	"strings"

	"ganjineh-auth/pkg/ierror"
	responses "ganjineh-auth/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ErrorHandler struct{}

func NewErrorHandler() *ErrorHandler {
    return &ErrorHandler{}
}
func (h *ErrorHandler) FiberErrorHandler(c *fiber.Ctx, err error) error {

    // AuthError (JWT/Blacklist/Refresh)
    if authErr, ok := err.(*ierror.AuthError); ok {

        // GraphQL
        if strings.HasPrefix(c.Path(), "/graphql") {
            return c.Status(authErr.HttpStatus).JSON(fiber.Map{
                "errors": []fiber.Map{
                    {
                        "message": authErr.Message,
                        "extensions": fiber.Map{
                            "code":   authErr.Code,
                            "reason": authErr.Reason,
                        },
                    },
                },
            })
        }

        // REST
        return c.Status(authErr.HttpStatus).JSON(responses.ErrorResponse(authErr))
    }

    // AppError (Business logic)
    if appErr, ok := err.(*ierror.AppError); ok {

        // GraphQL
        if strings.HasPrefix(c.Path(), "/graphql") {
            return c.Status(appErr.HttpStatus).JSON(fiber.Map{
                "errors": []fiber.Map{
                    {
                        "message": appErr.Message,
                        "extensions": fiber.Map{
                            "code": appErr.Code,
                        },
                    },
                },
            })
        }

        // REST
        return c.Status(appErr.HttpStatus).JSON(responses.ErrorResponse(appErr))
    }

    // Unknown error
    return c.Status(500).JSON(responses.BaseResponse{
        Code: 500,
        Error: &responses.ErrorInfo{
            Message: "internal server error",
        },
    })
}
