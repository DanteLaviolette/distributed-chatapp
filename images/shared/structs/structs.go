package structs

// Data type representing User stored in DB
type User struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}

// Data type representing user log in request
type LoginRequest struct {
	Email    string
	Password string
}

type RefreshTokenDocument struct {
	UserId string
	Secret string
}
