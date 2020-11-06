package test

import (
	"testing"

	"github.com/Improwised/golang-api/config"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	cfg := config.LoadTestEnv()

	assert.Equal(t, cfg.Env, "testing")

}
