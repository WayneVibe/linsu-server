package auths

import (
	"github.com/gin-gonic/gin"
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
func (h *Handler) Register(c *gin.Context) {
	var reqData RegisterReq
	if err := req.JsonParam(c, &reqData); err != nil {
		return
	}
	resp, err := h.service.register(reqData)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}
