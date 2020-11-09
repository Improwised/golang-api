package models

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/rs/xid"
)

// UserTable represent table name
const UserTable = "users"

// User model
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name" db:"first_name" validate:"required"`
	LastName  string `json:"last_name" db:"last_name" validate:"required"`
	Email     string `json:"email" db:"email" validate:"required"`
	CreatedAt string `json:"created_at" db:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at" db:"updated_at,omitempty"`
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
func (model *UserModel) GetUser() ([]User, error) {
	var users []User
	if err := model.db.From(UserTable).ScanStructs(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// InsertUser retrieve user
func (model *UserModel) InsertUser(user *User) error {
	ds := model.db.Insert(UserTable).Returning(goqu.C("id")).Rows(
		goqu.Record{
			"id":         xid.New().String(),
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		}).Executor()

	_, err := ds.ScanVal(&user.ID)
	if err != nil {
		return err
	}

	return nil
}
