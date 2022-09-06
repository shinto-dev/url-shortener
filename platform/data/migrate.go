package data

import (
	"context"
	"database/sql"
	"fmt"
	"url-shortner/platform/observation"
	"url-shortner/platform/observation/logging"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

// DBConfig defines the settings to connect to a database.
type DBConfig struct {
	Hostname string
	Port     int

	Database string
	Username string
	Password string
	DebugLog bool
}

// MigrateUp performs all pending migrations.
func MigrateUp(ctx context.Context, dbConfig DBConfig, migrationsFilePath string) error {
	observation.Add(ctx,
		logging.LField("context", "migration"),
		logging.LField("db_name", dbConfig.Database),
		logging.LField("db_host", dbConfig.Hostname),
	)

	m, err := getMigrate(dbConfig, migrationsFilePath)
	if err != nil {
		return errors.WithMessage(err, "getting migration object failed")
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			logging.FromContext(ctx).Info("No changes")
			return nil
		}
		return err
	}
	logging.FromContext(ctx).Info("migration successful")

	return nil
}

func getMigrate(dbConfig DBConfig, migrationsFilePath string) (*migrate.Migrate, error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=true",
		dbConfig.Username, dbConfig.Password, dbConfig.Hostname, dbConfig.Port, dbConfig.Database)

	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate database config. invalid DB url")
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return migrate.NewWithDatabaseInstance(migrationsFilePath, dbConfig.Database, driver)
}
