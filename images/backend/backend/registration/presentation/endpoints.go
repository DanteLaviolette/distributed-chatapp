package presentation

import (
	"net/http"

	"go.violettedev.com/eecs4222/registration/business"
	"go.violettedev.com/eecs4222/registration/structs"
	"go.violettedev.com/eecs4222/shared/auth"

	"github.com/gofiber/fiber/v2"
)

/*
Registration endpoint -- attempts to register user
Accepts a POST request containing RegistrationRequest as JSON.
Returns:
- 200 upon success
- 400 if any fields are empty or request is invalid
- 409 if email already exists
- 500 error code if unexpected error occurs
*/
func RegisterEndpoint(c *fiber.Ctx) error {
	// Parse body to struct
	var registerInfo structs.RegistrationRequest
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
func ChangePasswordEndpoint(c *fiber.Ctx) error {
	// Parse password from request
	var changePasswordInfo structs.ChangePasswordRequest
	if err := c.BodyParser(&changePasswordInfo); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	// Parse user id from locals (set by auth middleware)
	userId := c.Locals(auth.LOCALS_USER_ID)
	userIdString, ok := userId.(string)
	if ok {
		// Change password
		return c.SendStatus(
			business.ChangeUserPassword(userIdString, changePasswordInfo.Password),
		)
	} else {
		return c.SendStatus(http.StatusBadRequest)
	}
}
