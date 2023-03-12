package main

import (
	"log"
	"net/http"
	"os"

	"registration/business"
	"registration/structs"
	"shared/persistence"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

/*
Load dev environment from .env if os.Getenv('GO_ENV') != 'prod'
*/
func loadDevEnv() {
	// Get dev env if not prod
	if os.Getenv("GO_ENV") != "prod" {
		// Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

}

/*
Registration endpoint -- attempts to register user
Accepts a POST request containing RegisterInfo as JSON.
Returns:
- 400 if request is invalid (bad content type, bad method or invalid JSON)
-
*/
func registerEndpoint(c *fiber.Ctx) error {
	// Parse body to struct
	var registerInfo structs.RegisterInfo
	if err := c.BodyParser(&registerInfo); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	// Handle business logic
	res, resCode := business.RegisterUser(registerInfo)
	return c.Status(resCode).SendString(res)
}

func main() {
	loadDevEnv()
	persistence.InitializeDBConnection()
	app := fiber.New()
	app.Post("/api/register", registerEndpoint)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
