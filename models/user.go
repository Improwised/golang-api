package models

import (
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/database"
	"github.com/doug-martin/goqu/v9"
)

// UserTable represent table name
const UserTable = "users"

// User model
type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

var users []User

// UserModel implements user related database operations
type UserModel struct {
	db *goqu.Database
}

// InitUserModel Init model
func InitUserModel(cfg config.DBConfig) (*UserModel, error) {
	db, err := database.Connect(cfg)
	if err != nil {
		return nil, err
	}
	return &UserModel{
		db: db,
	}, nil
}

// GetUser retrive user
func (model *UserModel) GetUser() ([]User, error) {
	if err := model.db.From(UserTable).ScanStructs(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// InsertUser retrive user
func (model *UserModel) InsertUser(user *User) (int64, error) {
	preparedDs := model.db.From(UserTable).Prepared(true)
	ds, err := preparedDs.Insert().Rows(
		goqu.Record{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		}).Executor().Exec()

	id, err := ds.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
