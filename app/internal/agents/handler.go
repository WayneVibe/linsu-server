package agents

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/setcreed/hade-kit/req"
	"github.com/setcreed/hade-kit/res"
)

type Handler struct {
	service *service
}

func NewHandler() *Handler {
	return &Handler{
		service: newService(),
	}
}

func (h *Handler) CreateAgent(c *gin.Context) {
	var createReq CreateAgentReq
	if err := req.JsonParam(c, &createReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	// 如果需要做链路追踪 上下文要进行传递
	// 这个上下文超时是10s
	resp, err := h.service.createAgent(c.Request.Context(), userID, createReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func (h *Handler) ListAgents(c *gin.Context) {
	var listReq SearchAgentReq
	if err := req.JsonParam(c, &listReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	resp, err := h.service.listAgents(c.Request.Context(), userID, listReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func (h *Handler) GetAgent(c *gin.Context) {
	var id uuid.UUID
	if err := req.Path(c, "id", &id); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	resp, err := h.service.getAgent(c.Request.Context(), userID, id)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func (h *Handler) UpdateAgent(c *gin.Context) {
	var updateReq UpdateAgentReq
	if err := req.JsonParam(c, &updateReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	resp, err := h.service.updateAgent(c.Request.Context(), userID, updateReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}
