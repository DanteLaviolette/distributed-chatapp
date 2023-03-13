package business

import (
	"log"
	"net/http"
	"net/mail"
	"registration/persistence"
	"shared/structs"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const internalServerError = "Something went wrong. Try again later."
const missingFieldsError = "Fields can't be empty"
const passwordInvalidError = "Password invalid"
const emailInvalidError = "Email invalid"
const duplicateKeyError = "Email already exists"

/*
Registers the user if possible based on User. Returns the response
message and response code.
- 200 upon success
- 400 if any fields are empty or invalid
- 409 if email already exists
- 500 error code if unexpected error occurs
*/
func RegisterUser(user structs.User) (string, int) {
	// Fail if fields are missing
	if user.FirstName == "" || user.LastName == "" ||
		user.Email == "" || user.Password == "" {
		return missingFieldsError, http.StatusBadRequest
	}
	// Validate email
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return emailInvalidError, http.StatusBadRequest
	}
	// Validate password
	if !isPasswordValid(user.Password) {
		return passwordInvalidError, http.StatusBadRequest
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

/*
Changes the user password. Returns:
- 200 on success
- 400 on bad request (ie. invalid password)
- 500 if an error occurs
*/
func ChangeUserPassword(userId string, password string) int {
	// Validate password
	if !isPasswordValid(password) {
		return http.StatusBadRequest
	}
	// Get password hash
	hash, err := getPasswordHash(password)
	if err != nil {
		log.Print(err)
		return http.StatusInternalServerError
	}
	// Update users password in db
	err = persistence.UpdatePasswordForUserId(userId, hash)
	if err != nil {
		log.Print(err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// Returns true if the password is valid, false otherwise
func isPasswordValid(password string) bool {
	return len(password) >= 8
}

// Returns the hash of the password using bcrypt
func getPasswordHash(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(res), err
}
