package main

import (
	"log"
	"os"

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

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	routes.Router(app)

	app.Listen(host + ":" + port)

	defer database.DisconnectDB()
	defer database.CancelDBContext()
}
