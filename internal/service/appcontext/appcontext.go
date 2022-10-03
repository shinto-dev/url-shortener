package appcontext

import (
	"github.com/shinto-dev/url-shortener/config"
	"github.com/shinto-dev/url-shortener/foundation/data"
	"github.com/shinto-dev/url-shortener/internal/business/shorturl"
)

type AppContext struct {
	ShortURLCore shorturl.Core
}

func Get(conf config.Config) (AppContext, error) {
	db, err := data.Connect(data.DBConfig{
		Hostname: conf.Database.Hostname,
		Port:     conf.Database.Port,
		Database: conf.Database.DatabaseName,
		Username: conf.Database.Username,
		Password: conf.Database.Password,
		DebugLog: conf.Database.DebugLog,
	})
	if err != nil {
		return AppContext{}, err
	}

	return AppContext{
		ShortURLCore: shorturl.NewShortURLCore(db),
	}, nil
}
