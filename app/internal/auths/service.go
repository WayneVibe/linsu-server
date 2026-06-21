package auths

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"model"
	"net/smtp"
	"time"

	"github.com/google/uuid"
	"github.com/setcreed/hade-kit/cache"
	"github.com/setcreed/hade-kit/config"
	"github.com/setcreed/hade-kit/database"
	"github.com/setcreed/hade-kit/errs"
	"github.com/setcreed/hade-kit/logs"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"common/biz"
)

type service struct {
	repo  repository
	cache *cache.RedisCache
}

func newService() *service {
	return &service{
		repo:  newModel(database.GetPostgresDB().GormDB),
		cache: cache.NewRedisCache(),
	}
}

func (s *service) register(req RegisterReq) (any, error) {
	// 先检查用户名、邮箱是否已经注册
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Check if username already exists
	u, err := s.repo.findByUserName(ctx, req.Username)
	if err != nil {
		logs.Errorf("find auths err:%v", err)
		return nil, errs.DBError
	}
	if u != nil {
		return nil, biz.ErrUserNameExisted
	}

	// Check if email already exists
	u, err = s.repo.findByEmail(ctx, req.Email)
	if err != nil {
		logs.Errorf("find email err:%v", err)
		return nil, errs.DBError
	}
	if u != nil {
		return nil, biz.ErrEmailExisted
	}

	// 对密码进行加密
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, biz.ErrPasswordFormat
	}

	// 生成邮件激活的token
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, errs.DBError
	}
	verifyToken := hex.EncodeToString(tokenBytes)
	userId := uuid.New()

	// 存入redis中，用于激活邮件时验证
	redisCache := cache.NewRedisCache()
	tokenKey := fmt.Sprintf("verify_token:%s", verifyToken)
	logs.Infof("storing verify token in Redis: %s and userId: %s", tokenKey, userId)
	if err := redisCache.Set(tokenKey, userId.String(), 24*60*60); err != nil { // 24 hours
		logs.Errorf("failed to store verify token in Redis: %v", err)
		return nil, errs.DBError
	}

	u = &model.User{
		Id:            userId,
		Username:      req.Username,
		Password:      string(password),
		LastLoginTime: time.Now(),
		Status:        model.UserStatusPending, // Set status to pending until email verification
		Avatar:        "default",
		CurrentPlan:   model.FreePlan,
		Email:         req.Email,
		EmailVerified: false,
	}

	err = s.repo.transaction(func(tx *gorm.DB) error {
		if err := s.repo.saveUser(ctx, tx, u); err != nil {
			logs.Errorf("create user err:%v", err)
			return err
		}
		// 发送验证邮件
		if err := s.sendVerificationEmail(u.Email, u.Username, verifyToken); err != nil {
			logs.Errorf("send verification email err:%v", err)
			// Don't fail the registration if email sending fails, but log the error
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errs.DBError
	}
	return &RegisterResp{
		Message: "注册成功，请检查您的邮箱并点击验证链接完成注册",
	}, nil
}

func (s *service) sendVerificationEmail(email string, username string, token string) error {
	// 加载邮件配置
	emailConfig := config.GetConfig().Email

	if emailConfig.Host == nil || emailConfig.Port == nil {
		logs.Warn("Email not configured, skipping verification email")
		return nil
	}
	// Email content
	subject := "请验证您的邮箱地址"
	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", emailConfig.GetBaseURL(), token)
	body := fmt.Sprintf("尊敬的 %s，\n\n感谢您注册我们的服务！\n\n请点击以下链接验证您的邮箱地址：\n%s\n\n如果链接无法点击，"+
		"请复制并粘贴到浏览器地址栏中。\n\n谢谢！\n", username, verifyURL)

	// Set up authentication information
	auth := smtp.PlainAuth("", emailConfig.GetUsername(), emailConfig.GetPassword(), emailConfig.GetHost())

	// Connect to the server, authenticate, and send the email
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%d", emailConfig.GetHost(), emailConfig.GetPort())
	err := smtp.SendMail(addr, auth, emailConfig.GetFrom(), to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
