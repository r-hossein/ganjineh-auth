package models

type OTPPhoneRequest struct {
    PhoneNumber string `json:"phonenumber" form:"phonenumber" query:"phonenumber" validate:"required,len=11,startswith=09"`
}

type OTPVerifyRequest struct {
    PhoneNumber string `json:"phonenumber" validate:"required,len=11,startswith=09"`
    Signature string `json:"signature" validate:"required,len=16"`
    Code string `json:"code" validate:"required,len=6"`
}