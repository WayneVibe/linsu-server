package inits

import (
	"github.com/setcreed/hade-kit/config"
	"github.com/setcreed/hade-kit/database"
	"github.com/setcreed/hade-kit/server"
	"github.com/setcreed/hade-kit/tools/jwt"

	"app/internal/router"
)

func Init(s *server.Server, conf *config.Config) {
	// 初始化数据库
	database.InitPostgres(conf.DB.Postgres)
	// 初始化Redis
	database.InitRedis(conf.DB.Redis)
	// 初始化jwt
	jwt.Init(conf.Jwt.GetSecret())
	s.RegisterRouters(
		&router.Event{},
		&router.AuthRouter{},
		&router.AgentRouter{},
	)
}
