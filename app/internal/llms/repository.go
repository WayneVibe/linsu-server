package llms

import (
	"context"

	"github.com/google/uuid"

	"model"
)

type repository interface {
	createProviderConfig(ctx context.Context, m *model.ProviderConfig) error
	listProviderConfigs(ctx context.Context, userId uuid.UUID) ([]*model.ProviderConfig, int64, error)
	createLLM(ctx context.Context, llm *model.LLM) error
	listLLMs(ctx context.Context, userID uuid.UUID, filter LLMFilter) ([]*model.LLM, int64, error)
	listLLMAll(ctx context.Context, userID uuid.UUID, filter LLMFilter) ([]*model.LLM, error)
	getProviderConfig(ctx context.Context, provider string) (*model.ProviderConfig, error)
}
