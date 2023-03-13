package main

import (
	"log"
	"net/http"
	"os"

	"registration/business"
	"shared/auth"
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

/*
Changes the user password. Returns:
- 200 on success
- 400 on bad request (ie. invalid password or request)
- 401 if not signed in
- 500 if an error occurs
*/
func changePasswordEndpoint(c *fiber.Ctx) error {
	log.Print("hit")
	// Parse password from request
	var changePasswordInfo structs.ChangePasswordRequest
	if err := c.BodyParser(&changePasswordInfo); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	// Parse user id from locals (set by auth middleware)
	userId := c.Locals(auth.LOCALS_USER_ID)
	userIdString, ok := userId.(string)
	log.Print("hit")
	if ok {
		log.Print("hit 2")
		// Change password
		return c.SendStatus(
			business.ChangeUserPassword(userIdString, changePasswordInfo.Password),
		)
	} else {
		return c.SendStatus(http.StatusBadRequest)
	}
}

func main() {
	loadDevEnv()
	persistence.InitializeDBConnection(os.Getenv("MONGODB_URL"))
	authProvider := auth.Initialize(os.Getenv("AUTH_PRIVATE_KEY"),
		os.Getenv("REFRESH_PRIVATE_KEY"))
	app := fiber.New()
	app.Post("/api_register/register", registerEndpoint)
	app.Post("/api_register/change_password", authProvider.IsAuthenticatedFiberMiddleware, changePasswordEndpoint)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
