package pkg

import (
	"encoding/json"
	"strings"

	"ganjineh-auth/internal/repositories/db"
	"ganjineh-auth/internal/services"
	"ganjineh-auth/pkg/ierror"
	responses "ganjineh-auth/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ErrorHandler struct{
    backgroundService services.BackgroundServiceInterface
}

func NewErrorHandler(backgroundService services.BackgroundServiceInterface) *ErrorHandler {
    return &ErrorHandler{
        backgroundService: backgroundService,
    }
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
        
        requestBody, _ := h.getRequestBodyJSONB(c)
        //store unKnown error in database 
        h.backgroundService.InsertError(c.Context(),db.InsertErrorParams{
            HttpCode: int32(appErr.HttpStatus),
            StatusCode: int32(appErr.Code),
            Message: appErr.Message,
            StackTrace: &appErr.StackTrace,
            Endpoint: &c.Route().Path,
            Column6: requestBody,
        })

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

func (h *ErrorHandler) getRequestBodyJSONB(c *fiber.Ctx) ([]byte, error) {
    body := c.Body()
    if len(body) == 0 {
        return json.Marshal(map[string]interface{}{})
    }
    
    // بررسی اینکه آیا body معتبر JSON است
    var jsonBody interface{}
    if err := json.Unmarshal(body, &jsonBody); err != nil {
        // اگر JSON نیست، به عنوان string ذخیره کنید
        return json.Marshal(map[string]interface{}{
            "raw_body": string(body),
        })
    }
    
    return json.Marshal(jsonBody)
}