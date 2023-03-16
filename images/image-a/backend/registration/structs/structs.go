package structs

// Data type representing registration request to be created in DB
type RegistrationRequest struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}

// Data type representing change password request
type ChangePasswordRequest struct {
	Password string
}
