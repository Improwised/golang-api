package v1_test

import (
	"net/http"
	"testing"

	"github.com/Improwised/golang-api/pkg/structs"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	t.Run("create user invalid email", func(t *testing.T) {
		req := structs.ReqRegisterUser{
			Email: "abcd@@",
		}
		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Post("/api/v1/users")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})

	t.Run("create user with valid input", func(t *testing.T) {
		req := structs.ReqRegisterUser{
			Email:     "someone@example.com",
			FirstName: "someone",
			LastName:  "someone",
			Password:  "someone@1234",
			Roles:     "user",
		}
		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Post("/api/v1/users")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode())
	})

	t.Run("login user with invalid credentials", func(t *testing.T) {
		req := structs.ReqLoginUser{
			Email:    "someone@example.com",
			Password: "someonesomeone",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Post("/api/v1/login")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	})

	t.Cleanup(func() {
		_, err := db.Exec("delete from users where email='someone@example.com'")
		assert.Nil(t, err)
	})
}
