package v1_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Improwised/golang-api/constants"
	"github.com/Improwised/golang-api/pkg/structs"
	"github.com/stretchr/testify/assert"
)

func TestGetUserOnlyOwnRecord(t *testing.T) {
	register := func(email string) string {
		res, err := client.R().SetBody(structs.ReqRegisterUser{
			Email:     email,
			FirstName: "owner",
			LastName:  "check",
			Password:  "owner@1234",
			Roles:     "user",
		}).Post("/api/v1/users")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode())

		var body struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		assert.Nil(t, json.Unmarshal(res.Body(), &body))
		return body.Data.ID
	}

	aliceID := register("alice-owner@example.com")
	bobID := register("bob-owner@example.com")

	res, err := client.R().SetBody(structs.ReqLoginUser{
		Email:    "alice-owner@example.com",
		Password: "owner@1234",
	}).Post("/api/v1/login")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode())

	var aliceToken string
	for _, c := range res.Cookies() {
		if c.Name == constants.CookieUser {
			aliceToken = c.Value
		}
	}
	assert.NotEmpty(t, aliceToken)

	t.Run("own record", func(t *testing.T) {
		res, err := client.R().SetCookie(&http.Cookie{Name: constants.CookieUser, Value: aliceToken}).
			Get("/api/v1/users/" + aliceID)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	t.Run("another user's record", func(t *testing.T) {
		res, err := client.R().SetCookie(&http.Cookie{Name: constants.CookieUser, Value: aliceToken}).
			Get("/api/v1/users/" + bobID)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusForbidden, res.StatusCode())
	})

	t.Cleanup(func() {
		_, err := db.Exec("delete from users where email in ('alice-owner@example.com', 'bob-owner@example.com')")
		assert.Nil(t, err)
	})
}
