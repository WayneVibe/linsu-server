package inits

import (
	"github.com/setcreed/hade-kit/config"
	"github.com/setcreed/hade-kit/database"
	"github.com/setcreed/hade-kit/server"

	"app/internal/router"
)

func Init(s *server.Server, conf *config.Config) {
	database.InitPostgres(conf.DB.Postgres)
	database.InitRedis(conf.DB.Redis)
	s.RegisterRouters(
		&router.Event{},
		&router.AuthRouter{},
	)
}
