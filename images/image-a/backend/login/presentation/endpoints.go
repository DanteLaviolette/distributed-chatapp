package presentation

import (
	"net/http"

	"go.violettedev.com/eecs4222/auth"
	"go.violettedev.com/eecs4222/login/business"
	"go.violettedev.com/eecs4222/structs"

	"github.com/gofiber/fiber/v2"
)

/*
Login endpoint -- attempts to log user in, either returning an error
or an auth & refresh JWT cookie upon success. Accepts LoginRequest as JSON POST.
- Returns 200 upon success
- Returns 400 if given bad credentials
- Returns 500 if something else goes wrong
*/
func LoginEndpoint(c *fiber.Ctx) error {
	// Parse body to struct
	var loginRequest structs.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	// Return status code
	return c.SendStatus(business.Login(loginRequest, c))
}

/*
Logs the user out. Returns:
- 401 if not signed in
- 200 otherwise
*/
func LogoutEndpoint(c *fiber.Ctx) error {
	auth.InvalidateCredentials(c)
	return c.SendStatus(200)
}
