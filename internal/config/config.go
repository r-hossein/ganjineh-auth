package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/wire"
)

type StructConfig struct {
    AccessSecret      string        `mapstructure:"ACCESS_SECRET"`
    RefreshSecret     string        `mapstructure:"REFRESH_SECRET"`
    TempSecret        string        `mapstructure:"TEMP_SECRET"`
    AccessExpiration  time.Duration `mapstructure:"ACCESS_EXPIRATION"`
    RefreshExpiration time.Duration `mapstructure:"REFRESH_EXPIRATION"`
    TempExpiration    time.Duration `mapstructure:"TEMP_EXPIRATION"`

    BLUEPRINT_DB_HOST   string
    BLUEPRINT_DB_PORT   int
    BLUEPRINT_DB_DATABASE   string
    BLUEPRINT_DB_USERNAME   string
    BLUEPRINT_DB_PASSWORD   string
    BLUEPRINT_DB_SCHEMA string

    REDIS_HOST string
    REDIS_PORT int
    REDIS_USER string
    REDIS_PASS string

    PORT    int
    APP_ENV string

    SECRET_KEY string
}

func LoadConfig() *StructConfig {
    durationStr := getEnvWithDefault("JWT_ACCESS_EXPIRATION", "30m")
    accessExpiration, err := time.ParseDuration(durationStr)
    if err != nil {
        fmt.Printf("error in parsing JWT_ACCESS_EXPIRATION :%v\n",err)
        return nil
    }

    durationStr = getEnvWithDefault("JWT_REFRESH_EXPIRATION", "720h")
    refreshExpiration, err := time.ParseDuration(durationStr)
    if err != nil {
        fmt.Printf("error in parsing JWT_REFRESH_EXPIRATION :%v\n",err)
        return nil
    }
    durationStr = getEnvWithDefault("JWT_ACCESS_EXPIRATION", "30m")
    tempExpiration, err := time.ParseDuration(durationStr)
    if err != nil {
        fmt.Printf("error in parsing JWT_ACCESS_EXPIRATION :%v\n",err)
        return nil
    }

    port, _ := strconv.Atoi(os.Getenv("PORT"))
	portp, _ := strconv.Atoi(os.Getenv("BLUEPRINT_DB_PORT"))
    portr, _ := strconv.Atoi(os.Getenv(""))
    return &StructConfig{
        AccessSecret: os.Getenv("JWT_ACCESS_SECRET"),
        RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
        TempSecret: os.Getenv("JWT_TEMP_SECRET"),
        AccessExpiration: accessExpiration,
        RefreshExpiration: refreshExpiration,
        TempExpiration: tempExpiration,
        BLUEPRINT_DB_HOST: os.Getenv("BLUEPRINT_DB_HOST"),
        BLUEPRINT_DB_PORT: portp,
        BLUEPRINT_DB_DATABASE: os.Getenv("BLUEPRINT_DB_DATABASE"),
        BLUEPRINT_DB_USERNAME: os.Getenv("BLUEPRINT_DB_USERNAME"),
        BLUEPRINT_DB_PASSWORD: os.Getenv("BLUEPRINT_DB_PASSWORD"),
        BLUEPRINT_DB_SCHEMA: os.Getenv("BLUEPRINT_DB_SCHEMA"),
        REDIS_PORT: portr,
        REDIS_HOST: os.Getenv("REDIS_HOST"),
        REDIS_USER: os.Getenv("REDIS_USER"),
        REDIS_PASS: os.Getenv("REDIS_PASS"),
        PORT: port,
        APP_ENV: os.Getenv("APP_ENV"),
        SECRET_KEY: os.Getenv("SECRET_KEY"),
    }
}

var ConfigSet = wire.NewSet(LoadConfig)

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}