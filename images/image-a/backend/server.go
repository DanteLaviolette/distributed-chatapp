package main

import (
	"log"
	"os"

	"shared/auth"
	"shared/persistence"

	loginPresentation "login/presentation"
	registerPresentation "registration/presentation"

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

func main() {
	// Load dev environment if needed
	loadDevEnv()
	// Initialize db connection
	persistence.InitializeDBConnection(os.Getenv("MONGODB_URL"))
	// Create auth provider
	authProvider := auth.Initialize(os.Getenv("AUTH_PRIVATE_KEY"),
		os.Getenv("REFRESH_PRIVATE_KEY"))
	// Initialize fiber (REST framework)
	app := fiber.New()
	// Define endpoints
	app.Post("/api_login/login", loginPresentation.LoginEndpoint)
	app.Post("/api_login/logout", authProvider.IsAuthenticatedFiberMiddleware, loginPresentation.LogoutEndpoint)
	app.Post("/api_register/register", registerPresentation.RegisterEndpoint)
	app.Post("/api_register/change_password", authProvider.IsAuthenticatedFiberMiddleware, registerPresentation.ChangePasswordEndpoint)
	// Listen
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
