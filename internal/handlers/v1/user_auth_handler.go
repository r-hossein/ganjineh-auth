package v1

import (
	"ganjineh-auth/internal/services/auth"
	"ganjineh-auth/internal/models/requests"
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"ganjineh-auth/pkg/ierror"
    "ganjineh-auth/pkg/response"
)

type AuthHandlerInterface interface {
	RequestOTPHandler(c *fiber.Ctx) error
	VerifyOTPHandler(c *fiber.Ctx) error
}

type AuthHandler struct {
	AuthService auth.AuthService
	validate *validator.Validate
}

func (h *AuthHandler) RequestOTPHandler (c *fiber.Ctx) error {
	var req models.OTPPhoneRequest
	
	if err := c.BodyParser(&req); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse(ierror.ErrBadRequest))
	}

	// Validate input
    if err := h.validate.Struct(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse(ierror.ErrBadRequest))
    }

	ctx := c.Context()
	
	res,err := h.AuthService.RequestOTP(ctx, req.PhoneNumber)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse(err))
	}
	return c.JSON(responses.SuccessResponse(res,200))
}

func (h *AuthHandler) VerifyOTPHandler (c *fiber.Ctx) error {
	var req models.OTPVerifyRequest
	
	if err := c.BodyParser(&req); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse(ierror.ErrBadRequest))
	}

	// Validate input
    if err := h.validate.Struct(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse(ierror.ErrBadRequest))
    }

	ctx := c.Context()

	res,err := h.AuthService.VerifyOTP(ctx, &req)
	
	if err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(responses.ErrorResponse(err))
	}
	return c.JSON(responses.SuccessResponse(res,200))
}
