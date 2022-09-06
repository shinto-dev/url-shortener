package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("should load configs from env variable", func(t *testing.T) {
		os.Setenv("HTTPSERVER_PORT", "8083")
		os.Setenv("ENVIRONMENT", "test")

		defer func() {
			os.Unsetenv("HTTPSERVER_PORT")
			os.Unsetenv("ENVIRONMENT")
		}()

		conf := Config{}
		err := conf.Load()

		assert.NoError(t, err)
		assert.Equal(t, 8083, conf.HTTPServer.Port)
		assert.Equal(t, Environment("test"), conf.Environment)
	})
}
