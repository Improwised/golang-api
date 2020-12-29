package models

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/rs/xid"
	"gopkg.in/go-playground/validator.v9"
)

// UserTable represent table name
const UserTable = "users"

// User model
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name" db:"first_name" validate:"required"`
	LastName  string `json:"last_name" db:"last_name" validate:"required"`
	Email     string `json:"email" db:"email" validate:"required"`
	Password  string `json:"password,omitempty" db:"password" validate:"required"`
	Roles     string `json:"roles,omitempty" db:"roles" validate:"required"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// UserModel implements user related database operations
type UserModel struct {
	db *goqu.Database
}

// InitUserModel Init model
func InitUserModel(goqu *goqu.Database) (UserModel, error) {
	return UserModel{
		db: goqu,
	}, nil
}

// GetUser retrieve user
func (model *UserModel) GetUsers() ([]User, error) {
	var users []User
	if err := model.db.From(UserTable).ScanStructs(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUser retrieve user
func (model *UserModel) GetUser(user *User) error {
	if _, err := model.db.From(UserTable).Where(goqu.Ex{
		"id": user.ID,
	}).Select(
		"id",
		"first_name",
		"last_name",
		"email",
	).ScanStruct(user); err != nil {
		return err
	}
	return nil
}

// InsertUser retrieve user
func (model *UserModel) InsertUser(user *User) error {
	val := validator.New()

	user.ID = xid.New().String()

	err := val.Struct(user)
	if err != nil {
		return err
	}

	_, err = model.db.Insert(UserTable).Rows(
		goqu.Record{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"password":   user.Password,
			"roles":      user.Roles,
		},
	).Executor().Exec()
	if err != nil {
		return err
	}

	if err = model.GetUser(user); err != nil {
		return err
	}

	return nil
}
