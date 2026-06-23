package agents

import (
	"context"

	"github.com/google/uuid"

	"model"
)

type repository interface {
	createAgent(ctx context.Context, agent *model.Agent) error
	listAgents(ctx context.Context, userID uuid.UUID, filter AgentFilter) ([]*model.Agent, int64, error)
	getAgent(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Agent, error)
	updateAgent(ctx context.Context, agent *model.Agent) error
}
