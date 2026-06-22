package auths

type RegisterReq struct {
	Username string `json:"username" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
	Email    string `json:"email" binding:"required" validate:"required,email"`
}
type VerifyEmailReq struct {
	Token string `json:"token" form:"token" binding:"required" validate:"required"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
