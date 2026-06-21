package auths

import (
	"github.com/gin-gonic/gin"

	"github.com/setcreed/hade-kit/res"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}
func (h *Handler) Register(c *gin.Context) {
	res.Success(c, nil)
}
