package register

import (
	"fmt"
	"net/http"
	"structs"

	"golang.org/x/crypto/bcrypt"
)

/*
Registers the user if possible based on registerInfo. Returns the response
message and response code.
*/
func RegisterUser(registerInfo structs.RegisterInfo) (string, int) {
	// Fail if fields are missing
	if registerInfo.Username == "" || registerInfo.FirstName == "" ||
		registerInfo.LastName == "" || registerInfo.Email == "" ||
		registerInfo.Password == "" {
		return "Missing fields", http.StatusBadRequest
	}
	// Hash password
	hash, err := getPasswordHash(registerInfo.password)
	if err != nil {
		return "", http.StatusInternalServerError
	}
	fmt.Printf("%s", hash)
	// Attempt to register user
}

func getPasswordHash(password string) (string, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}
