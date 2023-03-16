package business

import (
	"log"
	"net/http"
	"os"

	"go.violettedev.com/eecs4222/login/persistence"
	"go.violettedev.com/eecs4222/login/structs"
	"go.violettedev.com/eecs4222/shared/auth"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

/*
Logs the user in setting auth header & refresh cookie upon success, returning a
a 200 status code.
Returns a failure status code if something goes wrong (400 for bad credentials,
or 500 for unexpected error)
*/
func Login(loginInfo structs.LoginRequest, c *fiber.Ctx) int {
	user, err := persistence.GetUser(loginInfo.Email)
	// User not found or network error case
	if err != nil {
		// Log network error
		if mongo.IsNetworkError(err) {
			log.Print(err)
		}
		return http.StatusBadRequest
	}
	// Invalid password case
	if !isPasswordEqual(loginInfo.Password, user.Password) {
		return http.StatusBadRequest
	}
	// Password is valid & user was found, generate auth & refresh token cookies
	authProvider := auth.Initialize(os.Getenv("AUTH_PRIVATE_KEY"),
		os.Getenv("REFRESH_PRIVATE_KEY"))
	// Handle credentials
	if !authProvider.GenerateAndSetCredentials(user, c) {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// Returns true if the password is equal to the hashed password. False otherwise
func isPasswordEqual(password string, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword),
		[]byte(password)) == nil
}
