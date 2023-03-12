package business

import (
	"net/http"
	"registration/structs"

	"golang.org/x/crypto/bcrypt"
)

/*
Registers the user if possible based on registerInfo. Returns the response
message and response code.
*/
func RegisterUser(registerInfo structs.RegisterInfo) (string, int) {
	// Fail if fields are missing
	if registerInfo.FirstName == "" || registerInfo.LastName == "" ||
		registerInfo.Email == "" || registerInfo.Password == "" {
		return "Missing fields", http.StatusBadRequest
	}
	// Hash password
	hash, err := getPasswordHash(registerInfo.Password)
	if err != nil {
		return "", http.StatusInternalServerError
	}
	// Attempt to register user

	return "", 200
}

func getPasswordHash(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(res), err
}
