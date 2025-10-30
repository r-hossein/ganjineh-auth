package server

import (
	"github.com/gofiber/fiber/v2"

	"ganjineh-auth/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "ganjineh-auth",
			AppName:      "ganjineh-auth",
		}),

		db: database.New(),
	}

	return server
}
