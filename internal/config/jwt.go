package config

import "time"

type JWTConfig struct {
    AccessSecret      string        `mapstructure:"ACCESS_SECRET"`
    RefreshSecret     string        `mapstructure:"REFRESH_SECRET"`
    TempSecret        string        `mapstructure:"TEMP_SECRET"`
    AccessExpiration  time.Duration `mapstructure:"ACCESS_EXPIRATION"`
    RefreshExpiration time.Duration `mapstructure:"REFRESH_EXPIRATION"`
    TempExpiration    time.Duration `mapstructure:"TEMP_EXPIRATION"`
}