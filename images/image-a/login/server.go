package main

import (
	"log"
	"net/http"
	"os"

	"login/business"
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
Login endpoint -- attempts to log user in, either returning an error
or an auth & refresh token upon success. Accepts LoginRequest as JSON POST.
*/
func loginEndpoint(c *fiber.Ctx) error {
	// Parse body to struct
	var loginRequest structs.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	// Handle business logic
	authCookie, refreshCookie, resCode := business.Login(loginRequest)
	// Set cookies if successful
	if resCode == 200 {
		c.Cookie(&authCookie)
		c.Cookie(&refreshCookie)
	}
	// Return status code
	return c.SendStatus(resCode)
}

func main() {
	loadDevEnv()
	persistence.InitializeDBConnection(os.Getenv("MONGODB_URL"))
	app := fiber.New()
	app.Post("/api/login", loginEndpoint)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
