package biz

import "github.com/setcreed/hade-kit/errs"

var (
	ErrUserNameExisted = errs.NewError(1001, "用户名已存在")
	ErrEmailExisted    = errs.NewError(1002, "邮箱已存在")
	ErrPasswordFormat  = errs.NewError(1003, "密码格式错误")
	ErrTokenInvalid    = errs.NewError(1004, "无效的token")
	ErrUserNotFound    = errs.NewError(1005, "用户不存在")
)
