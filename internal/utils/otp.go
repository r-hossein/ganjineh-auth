package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	// "encoding/base64"
	"encoding/hex"
	"time"

	"ganjineh-auth/internal/config"
	"ganjineh-auth/internal/models/entities"
	req "ganjineh-auth/internal/models/requests"
	"ganjineh-auth/pkg/ierror"

	"github.com/google/wire"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type OTPPkgInterface interface {
	GenerateOTP(phoneNumber string) (*models.OTP)
	VerifyOTP(data *models.OTP, reqData *req.OTPVerifyRequest) (bool, *ierror.AppError)
}

type OTPPkgStruct struct{
    Conf	*config.StructConfig
}

func NewOTPPkgService(cfg *config.StructConfig) *OTPPkgStruct {
    return &OTPPkgStruct{
        Conf: cfg,
    }
}

var OTPPkgSet = wire.NewSet(
	NewOTPPkgService,
	wire.Bind(new(OTPPkgInterface), new(*OTPPkgStruct)),
)

var _ OTPPkgInterface = (*OTPPkgStruct)(nil)

func (s *OTPPkgStruct) GenerateOTP(phoneNumber string) *models.OTP {
	
    // code,err := totp.GenerateCodeCustom(phoneNumber,time.Now(),totp.ValidateOpts{
    //     Period: 120,
    //     Digits: otp.DigitsSix,
    //     Algorithm: otp.AlgorithmSHA256,
    // })
    // if err != nil {
    //     return &models.OTP{
	// 	PhoneNumber : phoneNumber,
	// 	Code: "",
	// 	Signature: "",
	// 	ExpiresAt: time.Now(),
	// }
    // }
	code := "123456"
    signature := s.generateSignature(phoneNumber, code)

    result := &models.OTP{
		PhoneNumber : phoneNumber,
		Code: code,
		Signature: signature,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return result
}

func (s *OTPPkgStruct) VerifyOTP(data *models.OTP, reqData *req.OTPVerifyRequest) (bool, *ierror.AppError) {

    if time.Now().After(data.ExpiresAt) {
		return false,ierror.ErrOTPExpired
	}

	if !s.verifySignature(reqData) {
		return false, ierror.ErrInvalidOTP
	}

	// 3. بررسی کد OTP
	valid, err := totp.ValidateCustom(data.Code, data.PhoneNumber, time.Now(), totp.ValidateOpts{
		Period:    120,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA256,
	})
	if err != nil && valid {
		return false, ierror.ErrInvalidOTP
	}

	if data.Code != reqData.Code {
		return false, ierror.ErrInvalidOTP
	}
	return true,nil
}

// generateSignature تولید امضا با HMAC-SHA256
func (s *OTPPkgStruct) generateSignature(phoneNumber, code string) string {
	// // ایجاد داده برای امضا
	data := phoneNumber + "|" + code
	
	// ایجاد HMAC-SHA256
	h := hmac.New(sha256.New, []byte(s.Conf.SECRET_KEY))
	h.Write([]byte(data))
	
	// کدگذاری base64
	
	// return base64.URLEncoding.EncodeToString(h.Sum(nil))
	hexSignature := hex.EncodeToString(h.Sum(nil))
    return hexSignature[:16]
	//     data := phoneNumber + "|" + code
    
    // h := hmac.New(sha256.New, []byte(s.Conf.SECRET_KEY))
    // h.Write([]byte(data))
    
    // // Take first 16 bytes of the HMAC, then encode
    // hmacResult := h.Sum(nil)
    // return hex.EncodeToString(hmacResult[:16])
}

// verifySignature بررسی صحت امضا
func (s *OTPPkgStruct) verifySignature(reqData *req.OTPVerifyRequest) bool {
	expectedSignature := s.generateSignature(reqData.PhoneNumber, reqData.Code)
	return hmac.Equal([]byte(expectedSignature), []byte(reqData.Signature))
}