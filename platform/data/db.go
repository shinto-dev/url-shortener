package data

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connect returns a database connection using Gorm.
func Connect(config DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Hostname,
		config.Port,
		config.Database,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
