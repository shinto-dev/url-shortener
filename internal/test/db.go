package test

import (
	"testing"
	"url-shortener/platform/data"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func ConnectTestDB(t *testing.T) *gorm.DB {
	dbConfig := data.DBConfig{
		Hostname: "localhost",
		Port:     3306,
		Database: "short_url",
		Username: "root",
		Password: "root@123",
	}
	db, err := data.Connect(dbConfig)
	assert.NoError(t, err)

	return db
}
