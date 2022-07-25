package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/thenewsatria/seenaoo-backend/database"
)

func main() {
	// loading env variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading environtment variables")
	}

	database.ConnectDB(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"))
	database.PingDB()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")

	defer database.CancelDBContext()
	defer database.DisconnectDB()
}
