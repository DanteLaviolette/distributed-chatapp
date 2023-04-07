package business

import (
	"log"
	"net/http"
	"net/mail"

	"go.violettedev.com/eecs4222/registration/persistence"
	"go.violettedev.com/eecs4222/registration/structs"
	"go.violettedev.com/eecs4222/shared/database/schema"

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
func RegisterUser(registrationRequest structs.RegistrationRequest) (string, int) {
	valid, errMsg := validateRegistrationRequest(registrationRequest)
	if !valid {
		return errMsg, http.StatusBadRequest
	}
	// Hash password
	hashedPassword, err := getPasswordHash(registrationRequest.Password)
	if err != nil {
		log.Print(err)
		return internalServerError, http.StatusInternalServerError
	}
	// Create user
	user := schema.UserSchema{
		Email:     registrationRequest.Email,
		FirstName: registrationRequest.FirstName,
		LastName:  registrationRequest.LastName,
		Password:  hashedPassword,
	}
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

/*
Returns true, "" if the registration request is valid, false and an error string
otherwise.
*/
func validateRegistrationRequest(registrationRequest structs.RegistrationRequest) (bool, string) {
	// Fail if fields are missing
	if registrationRequest.FirstName == "" || registrationRequest.LastName == "" ||
		registrationRequest.Email == "" || registrationRequest.Password == "" {
		return false, missingFieldsError
	}
	// Validate email
	_, err := mail.ParseAddress(registrationRequest.Email)
	if err != nil {
		return false, emailInvalidError
	}
	// Validate password
	if !isPasswordValid(registrationRequest.Password) {
		return false, passwordInvalidError
	}
	return true, ""
}
