package agents

import (
	"context"
	"model"
	"time"

	"github.com/google/uuid"
	"github.com/setcreed/hade-kit/database"
	"github.com/setcreed/hade-kit/errs"
	"github.com/setcreed/hade-kit/logs"

	"common/biz"
)

type service struct {
	repo repository
}

func newService() *service {
	return &service{
		repo: newModel(database.GetPostgresDB().GormDB),
	}
}

func (s *service) createAgent(ctx context.Context, userId uuid.UUID, req CreateAgentReq) (any, error) {
	// 子上下文 不能超过10s
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	agent := model.DefaultAgent(userId, req.Name, req.Description, req.Status)
	err := s.repo.createAgent(ctx, agent)
	if err != nil {
		logs.Errorf("创建智能代理失败: %v", err)
		return nil, errs.DBError
	}
	return agent, nil
}

func (s *service) listAgents(ctx context.Context, userID uuid.UUID, req SearchAgentReq) (*ListAgentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	filter := AgentFilter{
		Name:   req.Params.Name,
		Status: req.Params.Status,
		Limit:  req.Params.PageSize,
		Offset: (req.Params.Page - 1) * req.Params.PageSize,
	}
	list, total, err := s.repo.listAgents(ctx, userID, filter)
	if err != nil {
		logs.Errorf("查询智能代理列表失败: %v", err)
		return nil, errs.DBError
	}
	return &ListAgentResponse{
		Agents: list,
		Total:  total,
	}, nil
}

func (s *service) getAgent(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Agent, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	agent, err := s.repo.getAgent(ctx, userID, id)
	if err != nil {
		logs.Errorf("查询智能代理失败: %v", err)
		return nil, errs.DBError
	}
	if agent == nil {
		return nil, biz.AgentNotFound
	}
	return agent, nil
}

func (s *service) updateAgent(ctx context.Context, userId uuid.UUID, req UpdateAgentReq) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	//先查询id是否存在
	agent, err := s.repo.getAgent(ctx, userId, req.ID)
	if err != nil {
		logs.Errorf("查询智能代理失败: %v", err)
		return nil, errs.DBError
	}
	if agent == nil {
		return nil, biz.AgentNotFound
	}

	//对更新的字段进行判断
	if req.Name != "" {
		agent.Name = req.Name
	}
	if req.Description != "" {
		agent.Description = req.Description
	}
	if req.Status != "" {
		agent.Status = req.Status
	}
	if req.SystemPrompt != "" {
		agent.SystemPrompt = req.SystemPrompt
	}
	if req.ModelProvider != "" {
		agent.ModelProvider = req.ModelProvider
	}
	if req.ModelName != "" {
		agent.ModelName = req.ModelName
	}
	if req.ModelParameters != nil {
		agent.ModelParameters = req.ModelParameters
	}
	if req.OpeningDialogue != "" {
		agent.OpeningDialogue = req.OpeningDialogue
	}
	err = s.repo.updateAgent(ctx, agent)
	if err != nil {
		logs.Errorf("更新智能代理失败: %v", err)
		return nil, errs.DBError
	}
	return agent, nil
}
