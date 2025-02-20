package server

import (
	"log"
	"os"

	"github.com/ednaldo-dilorenzo/iappointment/util"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func InitFiberApplication() *fiber.App {
	loadApplicationConfig()
	return fiber.New(fiber.Config{
		ErrorHandler: util.ServerErrorHandler,
	})
}

func loadApplicationConfig() {
	loadEnv()
	log.Println("Running: ", os.Getenv("CURRENT_ENV"))
}

func loadEnv() {
	env := os.Getenv("APPLICATION_ENV")

	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")

	if "test" != env {
		godotenv.Load(".env.local")
	}

	godotenv.Load(".env." + env)
	godotenv.Load()
}
