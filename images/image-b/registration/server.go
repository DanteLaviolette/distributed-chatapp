package main

import (
	"log"
	"net/http"
	"os"

	"registration/business"
	"shared/persistence"
	"shared/structs"

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
Accepts a POST request containing User as JSON.
Returns:
- 200 upon success
- 400 if any fields are empty or request is invalid
- 409 if email already exists
- 500 error code if unexpected error occurs
*/
func registerEndpoint(c *fiber.Ctx) error {
	// Parse body to struct
	var registerInfo structs.User
	if err := c.BodyParser(&registerInfo); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	// Handle business logic
	res, resCode := business.RegisterUser(registerInfo)
	return c.Status(resCode).SendString(res)
}

func main() {
	loadDevEnv()
	persistence.InitializeDBConnection(os.Getenv("MONGODB_URL"))
	app := fiber.New()
	app.Post("/api_register/register", registerEndpoint)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
