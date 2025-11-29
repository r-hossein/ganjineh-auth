package models

type OTPLoginResponse struct {
    PhoneNumber string `json:"phonenumber" form:"phonenumber" query:"phonenumber" validate:"required,len=11,startswith=09"`
    Signature string `json:"signture" validate:"required,len=16"`
}

type OTPVerifyResponse struct {
    AccessToken string `json:"access_token"`
    RefreshToken string `json:"refresh_token,omitempty"`
    ExpiresAt   int64   `json:"exp"`
    UserExists   bool   `json:"user_exists"`
    FirstName   string  `json:"first_name,omitempty"`
    LastName    string  `json:"last_name,omitempty"`
    UserID      string  `json:"user_id,omitempty"`
    Role        string  `json:"role,omitempty"`
    PhoneNumber string  `json:"phone_number"`
}

type RefreshToken struct {
    AccessToken string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresAt   int64   `json:"exp"`
}