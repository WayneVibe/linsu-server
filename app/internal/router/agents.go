package router

import (
	"app/internal/agents"

	"github.com/gin-gonic/gin"
)

type AgentRouter struct {
}

func (a *AgentRouter) Register(engine *gin.Engine) {
	agentsGroup := engine.Group("/api/v1/agents")
	{
		agentsHandler := agents.NewHandler()
		agentsGroup.POST("/create", agentsHandler.CreateAgent)
		agentsGroup.POST("/list", agentsHandler.ListAgents)
		agentsGroup.GET("/:id", agentsHandler.GetAgent)
		agentsGroup.PUT("/update", agentsHandler.UpdateAgent)
		agentsGroup.DELETE("/:id", agentsHandler.DeleteAgent)
		agentsGroup.POST("/chat", agentsHandler.AgentMessage)
	}

}
