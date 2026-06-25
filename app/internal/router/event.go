package router

import (
	"github.com/setcreed/hade-kit/event"

	"app/internal/llms"
)

type Event struct {
}

func (*Event) Register() {
	llmService := llms.NewPublicService()
	event.Register("getProviderConfig", llmService.GetProviderConfig)
}
