package database

import (
	"fmt"
	"sync"

	"github.com/Improwised/golang-api/models"
)

var (
	db []*models.User
	mu sync.Mutex
)

// Connect with database
func Connect() {
	db = make([]*models.User, 0)
	fmt.Println("Connected with Database")
}

func Insert(user *models.User) {
	mu.Lock()
	db = append(db, user)
	mu.Unlock()
}

func Get() []*models.User {
	return db
}
