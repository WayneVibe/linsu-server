package inits

import (
	"github.com/setcreed/hade-kit/config"
	"github.com/setcreed/hade-kit/server"

	"app/internal/router"
)

func Init(s *server.Server, conf *config.Config) {
	s.RegisterRouters(
		&router.Event{},
		&router.AuthRouter{},
	)
}
