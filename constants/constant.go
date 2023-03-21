package constants

// variables
const (
	CookieUser = "user"
)

// fiber contexts
const (
	ContextUid = "userId"
)

// params
const (
	ParamUid = "userId"
)

// Success messages
// ...

// Fail messages
// ...
const (
	Unauthenticated    = "unauthenticated to access resource"
	Unauthorized       = "unauthorized to access resource"
	InvalidCredentials = "invalid credenticals"
	UserNotExist       = "user does not exists"
)

// Error messages
const (
	ErrGetUser         = "error while get user"
	ErrLoginUser       = "error while login user"
	ErrInsertUser      = "error while creating user, please try after sometime"
	ErrHealthCheckDb   = "error while checking health of database"
	ErrUnauthenticated = "error verifing user identity"
)

// Events
const (
	EventUserRegistered = "event:userRegistered"
)
