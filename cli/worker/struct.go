package worker

type UpdateUser struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Roles     string `json:"roles"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type DeleteUser struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Roles     string `json:"roles"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type AddUser struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Roles     string `json:"roles"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

