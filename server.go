package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/thenewsatria/seenaoo-backend/api/routes"
	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/database/seeds"
)

func main() {
	// loading env variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading environtment variables")
	}

	database.ConnectDB()
	database.PingDB()
	seeds.SeedPermissionsCollection()

	app := fiber.New()

	routes.Router(app)

	app.Listen(":3000")

	defer database.DisconnectDB()
	defer database.CancelDBContext()
}
