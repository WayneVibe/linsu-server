package auths

type RegisterReq struct {
	Username string `json:"username" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
	Email    string `json:"email" binding:"required" validate:"required,email"`
}
