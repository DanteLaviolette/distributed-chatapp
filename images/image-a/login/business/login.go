package business

import (
	"log"
	"login/persistence"
	"net/http"
	"os"
	"shared/auth"
	"shared/structs"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

/*
Logs the user in returning auth & refresh JWT cookies upon success along with
a 200 status code.
Returns a failure status code if something goes wrong (400 for bad credentials,
or 500 for unexpected error)
*/
func Login(loginInfo structs.LoginRequest) (*fiber.Cookie, *fiber.Cookie, int) {
	user, err := persistence.GetUserWithId(loginInfo.Email)
	// User not found or network error case
	if err != nil {
		// Log network error
		if mongo.IsNetworkError(err) {
			log.Print(err)
		}
		return nil, nil, http.StatusBadRequest
	}
	// Invalid password case
	if !isPasswordEqual(loginInfo.Password, user.Password) {
		return nil, nil, http.StatusBadRequest
	}
	// Password is valid & user was found, generate auth & refresh token cookies
	authProvider := auth.Initialize(os.Getenv("AUTH_PRIVATE_KEY"),
		os.Getenv("REFRESH_PRIVATE_KEY"))
	authToken := authProvider.CreateAuthCookie(user)
	refreshToken := authProvider.CreateRefreshCookie(user)
	if authToken == nil || refreshToken == nil {
		return nil, nil, http.StatusInternalServerError
	}
	return authToken, refreshToken, http.StatusOK
}

// Returns true if the password is equal to the hashed password. False otherwise
func isPasswordEqual(password string, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword),
		[]byte(password)) == nil
}
