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

func (h *Handler) VerifyEmail(c *gin.Context) {
	var reqData VerifyEmailReq
	if err := req.QueryParam(c, &reqData); err != nil {
		return
	}

	_, err := h.service.verifyEmail(reqData.Token)
	if err != nil {
		res.Error(c, err)
		return
	}
	//如果成功 直接跳转登录页面
	c.Redirect(302, "http://localhost:5173/login")
}

func (h *Handler) Login(c *gin.Context) {
	var loginReq LoginReq
	if err := req.JsonParam(c, &loginReq); err != nil {
		return
	}
	resp, err := h.service.login(loginReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}
