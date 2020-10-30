package models

import (
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/database"
	"github.com/doug-martin/goqu/v9"
)

// User model
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users []User

// UserModel implements user related database operations
type UserModel struct {
	db *goqu.Database
}

// InitUserModel Init model
func InitUserModel(config config.DBConfig) (*UserModel, error) {
	db, err := database.Connect(config)
	if err != nil {
		return nil, err
	}
	return &UserModel{
		db: db,
	}, nil
}

// GetUser retrive user
func (model *UserModel) GetUser() ([]User, error) {
	if err := model.db.From("user").ScanStructs(&users); err != nil {
		return nil, err
	}
	return users, nil
}
