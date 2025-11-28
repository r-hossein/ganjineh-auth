package models

import "time"

type SessionType string

const (
    SessionTypeActive  SessionType = "active"
    SessionTypeRevoke  SessionType = "revoked"
    SessionTypeUpdate  SessionType = "updated"
)

type BlackListSession struct {
	SessionId string `json:"sid"`
	Status SessionType		`json:"status"`
	ExpiresAt time.Time `json:"expiresat"`
}