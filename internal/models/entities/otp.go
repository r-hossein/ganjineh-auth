package models

import "time"

type OTP struct {
	PhoneNumber string `json:"phonenumber"`
	Code string		`json:"code"`
	Signature string `json:"signature"`
	ExpiresAt time.Time `json:"expiresat"`
}