package business

import (
	"log"
	"net/http"
	"registration/persistence"
	"shared/structs"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const internalServerError = "Something went wrong. Try again later."
const missingFieldsError = "Fields can't be empty"
const duplicateKeyError = "Email already exists"

/*
Registers the user if possible based on User. Returns the response
message and response code.
- 200 upon success
- 400 if any fields are empty
- 409 if email already exists
- 500 error code if unexpected error occurs
*/
func RegisterUser(user structs.User) (string, int) {
	// Fail if fields are missing
	if user.FirstName == "" || user.LastName == "" ||
		user.Email == "" || user.Password == "" {
		return missingFieldsError, http.StatusBadRequest
	}
	// Hash password
	hash, err := getPasswordHash(user.Password)
	if err != nil {
		log.Print(err)
		return internalServerError, http.StatusInternalServerError
	}
	// Set password to the hashed password
	user.Password = hash
	// Attempt to register user
	err = persistence.InsertUser(user)
	// Fail if error ocurred (dupe key error means email already exists)
	if mongo.IsDuplicateKeyError(err) {
		return duplicateKeyError, http.StatusConflict
	} else if err != nil {
		log.Print(err)
		return internalServerError, http.StatusInternalServerError
	}
	return "", 200
}

func getPasswordHash(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(res), err
}
