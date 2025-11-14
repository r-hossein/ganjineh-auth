package otp

import (
	"context"
	"ganjineh-auth/internal/repositories"
	"ganjineh-auth/internal/utils"
	req "ganjineh-auth/internal/models/requests"
	res "ganjineh-auth/internal/models/responses"
	"ganjineh-auth/pkg/ierror"
)

type OTPService interface {
	OTPRequest(ctx context.Context, phoneNumber string) (*res.OTPLoginResponse,*ierror.AppError)	
	ValidateOTP(ctx context.Context, data *req.OTPVerifyRequest ) (bool,*ierror.AppError)
}

type OTPResult struct{
	Signature string
}

type oTPService struct {
	otpRepo		repositories.RedisRepository 
	otpUtil 	utils.OTPInterface
}


func (o *oTPService) OTPRequest(ctx context.Context, phoneNumber string) (*res.OTPLoginResponse,*ierror.AppError){
	
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

func (o *oTPService) ValidateOTP(ctx context.Context, data *req.OTPVerifyRequest) (bool,*ierror.AppError) {
	
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
