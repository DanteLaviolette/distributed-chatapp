package business

import (
	"net/http"
	"registration/persistence"
	"registration/structs"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const internalServerError = "Something went wrong. Try again later."
const missingFieldsError = "Fields can't be empty"
const duplicateKeyError = "Email already exists"

/*
Registers the user if possible based on registerInfo. Returns the response
message and response code.
*/
func RegisterUser(registerInfo structs.RegisterInfo) (string, int) {
	// Fail if fields are missing
	if registerInfo.FirstName == "" || registerInfo.LastName == "" ||
		registerInfo.Email == "" || registerInfo.Password == "" {
		return missingFieldsError, http.StatusBadRequest
	}
	// Hash password
	hash, err := getPasswordHash(registerInfo.Password)
	if err != nil {
		return internalServerError, http.StatusInternalServerError
	}
	// Set password to the hashed password
	registerInfo.Password = hash
	// Attempt to register user
	err = persistence.InsertUser(registerInfo)
	// Fail if error ocurred (dupe key error means email already exists)
	if mongo.IsDuplicateKeyError(err) {
		return duplicateKeyError, http.StatusConflict
	} else if err != nil {
		return internalServerError, http.StatusInternalServerError
	}
	return "", 200
}

func getPasswordHash(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(res), err
}
