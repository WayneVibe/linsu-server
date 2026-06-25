package llms

import (
	"context"
	"time"

	"github.com/setcreed/hade-kit/database"
	"github.com/setcreed/hade-kit/event"
	"github.com/setcreed/hade-kit/logs"

	"app/shared"
)

type PublicService struct {
	repo repository
}

func NewPublicService() *PublicService {
	return &PublicService{
		repo: newModels(database.GetPostgresDB().GormDB),
	}
}

func (s *PublicService) GetProviderConfig(e event.Event) (any, error) {
	request := e.Data.(*shared.GetProviderConfigsRequest)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	providerConfig, err := s.repo.getProviderConfig(ctx, request.Provider)
	if err != nil {
		logs.Errorf("get provider config error: %v", err)
		return nil, err
	}
	return providerConfig, nil
}
