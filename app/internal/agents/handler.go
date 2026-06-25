package agents

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/setcreed/hade-kit/logs"
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

func (h *Handler) DeleteAgent(c *gin.Context) {
	var agentId uuid.UUID
	if err := req.Path(c, "id", &agentId); err != nil {
		return
	}
	err := h.service.deleteAgent(c.Request.Context(), agentId)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, nil)
}

func (h *Handler) AgentMessage(c *gin.Context) {
	//获取参数
	var messageReq AgentMessageReq
	if err := req.JsonParam(c, &messageReq); err != nil {
		return
	}
	userID, exist := req.GetUserIdUUID(c)
	if !exist {
		return
	}
	rc := http.NewResponseController(c.Writer)
	// 将当前请求的写入超时设置为零值（即无限制）
	// 这会覆盖全局 http.Server 的 WriteTimeout 设置
	if err := rc.SetWriteDeadline(time.Time{}); err != nil {
		// 如果失败，记录日志，但通常不会失败
		logs.Warn("Failed to set write deadline", "err", err)
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	// 这里我们用一个可以取消的context
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	// 这个接口是AI回答，我们返回两个chan，一个datachan 一个errchan
	// 调用大模型 我们需要放在协程处理，所以这里用channel
	dataChan, errorChan := h.service.agentMessage(ctx, userID, messageReq)
	//创建一个心跳 这里是防止一些防火墙拦截 导致连接中断
	heartbeat := time.NewTicker(time.Second * 5)
	defer heartbeat.Stop()

	for {
		//处理数据
		select {
		case <-ctx.Done():
			logs.Warnf("context done, 客户端断开连接")
			return
		case <-heartbeat.C:
			//处理心跳 我们发送一个冒号开头的消息 表示这是一个心跳消息
			_, err := c.Writer.Write([]byte(": keep-alive\n\n"))
			if err != nil {
				logs.Warnf("write heartbeat error: %v", err)
				cancel()
				return
			}
			//在go中处理消息 如果想要立即发送给客户端需要调用Flush
			c.Writer.Flush()
		case data, ok := <-dataChan:
			if !ok {
				//这里代表channel被关闭了 也就是消息结束了
				//按照SSE的规范，发送一个结束消息 [DONE]
				_, err := c.Writer.Write([]byte("data: [DONE]\n\n"))
				if err != nil {
					logs.Warnf("write done error: %v", err)
				}
				c.Writer.Flush()
				return
			}
			//有消息就直接发送， 这里我们不区分event 都按照默认message进行处理，前端也是如此
			//data数据是json的格式
			_, err := c.Writer.Write([]byte("data: " + data + "\n\n"))
			if err != nil {
				logs.Errorf("write data error: %v", err)
				cancel()
				return
			}
			c.Writer.Flush()
		case err, ok := <-errorChan:
			if !ok {
				//error的消息结束不处理，交给dataChan处理
				errorChan = nil
				continue
			}
			//如果有错误 发送错误的消息提供给客户端
			if err != nil {
				_, err := c.Writer.Write([]byte("data: [ERROR]" + err.Error() + "\n\n"))
				if err != nil {
					logs.Errorf("write error error: %v", err)
					cancel()
					return
				}
				c.Writer.Flush()
				return
			}
		}
	}
}
