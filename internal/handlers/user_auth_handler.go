package handlers

import (
	"ganjineh-auth/internal/models/requests"
	"ganjineh-auth/internal/services"
	"ganjineh-auth/pkg/ierror"
	"ganjineh-auth/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

type AuthHandlerInterface interface {
	RequestOTPHandler(c *fiber.Ctx) error
	VerifyOTPHandler(c *fiber.Ctx) error
	RegisterUserHandler(c *fiber.Ctx) error
}

type AuthHandlerStruct struct {
	AuthService services.AuthServiceInterface
	Validate *validator.Validate
}

func NewAuthHandler(authServ services.AuthServiceInterface, valid *validator.Validate) *AuthHandlerStruct{
	return &AuthHandlerStruct{
		AuthService: authServ,
		Validate: valid,
	}
}

var AuthHandlerSet = wire.NewSet(
	NewAuthHandler,
	wire.Bind(new(AuthHandlerInterface), new(*AuthHandlerStruct)),
)

func (h *AuthHandlerStruct) RequestOTPHandler (c *fiber.Ctx) error {
	var req models.OTPPhoneRequest
	
	if err := c.BodyParser(&req); err != nil{
		return ierror.ErrBadRequest
	}

	// Validate input
    if err := h.Validate.Struct(req); err != nil {
        return ierror.ErrBadRequest
    }

	ctx := c.Context()
	
	res,err := h.AuthService.RequestOTP(ctx, req.PhoneNumber)
	if err != nil {
		return err
	}
	return c.JSON(responses.SuccessResponse(res,200))
}

func (h *AuthHandlerStruct) VerifyOTPHandler (c *fiber.Ctx) error {
	var req models.OTPVerifyRequest
	
	if err := c.BodyParser(&req); err != nil{
		return ierror.ErrBadRequest
	}

	// Validate input
    if err := h.Validate.Struct(req); err != nil {
        return ierror.ErrBadRequest
    }

	ctx := c.Context()

	res,err := h.AuthService.VerifyOTP(ctx, &req)
	
	if err != nil {
		return err
	}
	return c.JSON(responses.SuccessResponse(res,200))
}

func (h *AuthHandlerStruct) RegisterUserHandler (c *fiber.Ctx) error {
	var req models.RegisterUserFirstRequest
	
	if err := c.BodyParser(&req); err != nil{
		return ierror.ErrBadRequest
	}

	// Validate input
    if err := h.Validate.Struct(req); err != nil {
        return ierror.ErrBadRequest
    }

	ctx := c.Context()

	res,err := h.AuthService.Register(ctx, &req)
	
	if err != nil {
		return err
	}
	return c.JSON(responses.SuccessResponse(res,200))
}
