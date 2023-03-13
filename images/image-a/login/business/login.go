package business

import (
	"log"
	"login/persistence"
	"net/http"
	"os"
	"shared/auth"
	"shared/structs"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(loginInfo structs.LoginRequest) (*fiber.Cookie, *fiber.Cookie, int) {
	user, err := persistence.GetUserWithId(loginInfo.Email)
	log.Print(user)
	// User not found case
	if err != nil {
		return nil, nil, http.StatusBadRequest
	}
	// Invalid password case
	if !isPasswordEqual(loginInfo.Password, user.Password) {
		return nil, nil, http.StatusBadRequest
	}
	// Password is valid & user was found, generate auth & refresh token
	authToken := auth.CreateAuthCookie(user, os.Getenv("AUTH_PRIVATE_KEY"))
	refreshToken := auth.CreateRefreshCookie(user, os.Getenv("REFRESH_PRIVATE_KEY"))
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
