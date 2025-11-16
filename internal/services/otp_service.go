package services

import (
	"context"
	req "ganjineh-auth/internal/models/requests"
	res "ganjineh-auth/internal/models/responses"
	"ganjineh-auth/internal/repositories"
	"ganjineh-auth/internal/utils"
	"ganjineh-auth/pkg/ierror"

	"github.com/google/wire"
)

type OTPServiceInterface interface {
	OTPRequest(ctx context.Context, phoneNumber string) (*res.OTPLoginResponse,*ierror.AppError)	
	ValidateOTP(ctx context.Context, data *req.OTPVerifyRequest ) (bool,*ierror.AppError)
}

type OTPResult struct{
	Signature string
}

type OTPServiceStruct struct {
	otpRepo		repositories.RedisOTPRepositoryInterface 
	otpUtil 	utils.OTPPkgInterface
}

func NewOTPService(
    otpRepo repositories.RedisOTPRepositoryInterface,
    otpUtil utils.OTPPkgInterface,
) OTPServiceInterface {
    return &OTPServiceStruct{
        otpRepo: otpRepo,
        otpUtil: otpUtil,
    }
}

var OTPServiceSet = wire.NewSet(
	NewOTPService,
	// wire.Bind(new(OTPServiceInterface), new(*OTPServiceStruct)),
)

var _ OTPServiceInterface = (*OTPServiceStruct)(nil)

func (o *OTPServiceStruct) OTPRequest(ctx context.Context, phoneNumber string) (*res.OTPLoginResponse,*ierror.AppError){
	
	result := o.otpUtil.GenerateOTP(phoneNumber)
	err := o.otpRepo.StoreOTP(ctx,result,5*60*1000)
	
	if err != nil {
		return nil, err
	}

	test := &res.OTPLoginResponse{
		PhoneNumber: phoneNumber,
		Signature: result.Signature,
	}
	return test,nil
}

func (o *OTPServiceStruct) ValidateOTP(ctx context.Context, data *req.OTPVerifyRequest) (bool,*ierror.AppError) {
	
	result,err := o.otpRepo.GetOTP(ctx, data.PhoneNumber)

	if err!= nil {
		return false,err
	}

	VerifyResult, err := o.otpUtil.VerifyOTP(result,data)

	if err != nil {
		return false, err
	}

	return VerifyResult,nil
}


func test() {
	panic("sajkdsf")
}