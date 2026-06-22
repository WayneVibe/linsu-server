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

	u, err := s.repo.findByUserName(ctx, req.Username)
	if err != nil {
		logs.Errorf("find auths err:%v", err)
		return nil, errs.DBError
	}
	if u != nil {
		return nil, biz.ErrUserNameExisted
	}

	// 检查邮箱是否已经注册
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
		logs.Errorf("register GenerateFromPassword error: %v", err)
		return nil, biz.ErrPasswordFormat
	}

	// 生成邮件用的token 用于邮件激活
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
		Status:        model.UserStatusPending,
		Avatar:        "default",
		CurrentPlan:   model.FreePlan,
		Email:         req.Email,
		EmailVerified: false,
	}

	err = s.repo.transaction(ctx, func(tx *gorm.DB) error {
		// 创建用户
		if err := s.repo.saveUser(ctx, tx, u); err != nil {
			logs.Errorf("create user err:%v", err)
			return err
		}
		// 发送验证邮件
		if err := s.sendVerificationEmail(u.Email, u.Username, verifyToken); err != nil {
			logs.Errorf("send verification email err:%v", err)
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
	//加载邮件的配置
	emailConfig := config.GetConfig().Email
	addr := fmt.Sprintf("%s:%d", emailConfig.GetHost(), emailConfig.GetPort())
	auth := smtp.PlainAuth("", emailConfig.GetUsername(), emailConfig.GetPassword(), emailConfig.GetHost())
	to := []string{email}
	subject := "请验证您的邮箱地址"
	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", emailConfig.GetBaseURL(), token)
	body := fmt.Sprintf("尊敬的 %s，\n\n感谢您注册我们的服务！\n\n请点击以下链接验证您的邮箱地址：\n%s\n\n如果链接无法点击，"+
		"请复制并粘贴到浏览器地址栏中。\n\n谢谢！\n", username, verifyURL)
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	err := smtp.SendMail(addr, auth, emailConfig.GetFrom(), to, msg)
	return err
}

func (s *service) verifyEmail(token string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 从Redis中获取用户ID
	tokenKey := fmt.Sprintf("verify_token:%s", token)
	userIdStr, err := s.cache.Get(tokenKey)
	if err != nil {
		logs.Errorf("verifyEmail Get error: %v", err)
		return nil, biz.ErrTokenInvalid
	}
	// 这个验证邮件的时候 redis的key也需要删除
	defer func() {
		err = s.cache.Set(tokenKey, "", 1)
		if err != nil {
			logs.Errorf("set email verify token cache error: %v", err)
		}
	}()

	//转换成uuid
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		logs.Errorf("verifyEmail Parse error: %v", err)
		return nil, biz.ErrTokenInvalid
	}

	//根据用户id查找 用户
	u, err := s.repo.findById(ctx, userId)
	if err != nil {
		logs.Errorf("verifyEmail findById error: %v", err)
		return nil, errs.DBError
	}
	if u == nil {
		return nil, biz.ErrUserNotFound
	}
	// 判断用户邮箱是否已经验证
	if u.EmailVerified {
		//直接返回验证成功
		return nil, nil
	}
	// 更新 用户
	u.EmailVerified = true
	u.Status = model.UserStatusNormal
	err = s.repo.transaction(ctx, func(tx *gorm.DB) error {
		return s.repo.updateUser(ctx, tx, u)

	})
	if err != nil {
		logs.Errorf("verifyEmail updateUser error: %v", err)
		return nil, errs.DBError
	}
	return nil, nil
}
