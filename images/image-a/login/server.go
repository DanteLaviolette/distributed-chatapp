package main

import (
	"log"
	"net/http"
	"os"

	"login/business"
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
Login endpoint -- attempts to log user in, either returning an error
or an auth & refresh JWT cookie upon success. Accepts LoginRequest as JSON POST.
- Returns 200 upon success
- Returns 400 if given bad credentials
- Returns 500 if something else goes wrong
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
		c.Cookie(authCookie)
		c.Cookie(refreshCookie)
	}
	// Return status code
	return c.SendStatus(resCode)
}

/*
Logs the user out. Returns:
- 401 if not signed in
- 200 otherwise
*/
func logoutEndpoint(c *fiber.Ctx) error {
	c.ClearCookie(auth.AUTH_COOKIE_NAME)
	c.ClearCookie(auth.REFRESH_COOKIE_NAME)
	c.SendStatus(200)
	// Get refresh token id (set by auth middleware)
	refreshTokenId := c.Locals(auth.LOCALS_REFRESH_TOKEN_ID)
	refreshTokenIdString, ok := refreshTokenId.(string)
	if ok {
		business.Logout(refreshTokenIdString)
	}
	return nil
}

func main() {
	loadDevEnv()
	persistence.InitializeDBConnection(os.Getenv("MONGODB_URL"))
	authProvider := auth.Initialize(os.Getenv("AUTH_PRIVATE_KEY"),
		os.Getenv("REFRESH_PRIVATE_KEY"))
	app := fiber.New()
	app.Post("/api_login/login", loginEndpoint)
	app.Post("/api_login/logout", authProvider.IsAuthenticatedFiberMiddleware, logoutEndpoint)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
