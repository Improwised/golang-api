package models

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/rs/xid"
)

// UserTable represent table name
const UserTable = "users"

// User model
type User struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	CreatedAt string `db:"created_at,omitempty"`
	UpdatedAt string `db:"updated_at,omitempty"`
}

// UserModel implements user related database operations
type UserModel struct {
	db *goqu.Database
}

// InitUserModel Init model
func InitUserModel(goqu *goqu.Database) (*UserModel, error) {
	return &UserModel{
		db: goqu,
	}, nil
}

// GetUser retrive user
func (model *UserModel) GetUser() ([]User, error) {
	var users []User
	if err := model.db.From(UserTable).ScanStructs(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// InsertUser retrive user
func (model *UserModel) InsertUser(user *User) (string, error) {
	var id string
	ds := model.db.Insert(UserTable).Returning(goqu.C("id")).Rows(
		goqu.Record{
			"id":         xid.New().String(),
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		}).Executor()

	_, err := ds.ScanVal(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
