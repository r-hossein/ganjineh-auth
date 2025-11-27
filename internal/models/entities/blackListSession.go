package models

import "time"

type BlackListSession struct {
	SessionId string `json:"sid"`
	Status string		`json:"status"`
	ExpiresAt time.Time `json:"expiresat"`
}