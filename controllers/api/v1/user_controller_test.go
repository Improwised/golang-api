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
}
