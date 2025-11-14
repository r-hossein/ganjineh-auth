package server

import (
	"github.com/gofiber/fiber/v2"
	"ganjineh-auth/internal/database"
	"ganjineh-auth/internal/repositories"
)

type FiberServer struct {
	*fiber.App
	pdb database.ServiceP
	rdb database.ServiceR
	repos *repositories.Container
}

func New() *FiberServer {
	pdb := database.NewP()
	rdb := database.NewR()
	repos := repositories.NewContainer(pdb, rdb)

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "ganjineh-auth",
			AppName:      "ganjineh-auth",
		}),
		pdb:   pdb,
		rdb:   rdb,
		repos: repos,
	}
	server.RegisterFiberRoutes()
	return server
}
