package jwt

import (
	"os"
	"testing"
	"time"

	"github.com/Improwised/golang-api/config"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	os.Chdir("../../")

	var token string = ""
	var err error = nil
	var subject string = "11112"
	cfg := config.LoadTestEnv()

	t.Run("create token", func(t *testing.T) {
		token, err = CreateToken(cfg, subject, time.Now().Add(time.Hour*1))
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("parse token", func(t *testing.T) {
		claims, err := ParseToken(cfg, token)
		assert.Nil(t, err)
		assert.Equal(t, subject, claims.Subject())
	})
}
