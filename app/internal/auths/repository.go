package auths

import (
	"context"

	"gorm.io/gorm"

	"model"
)

type repository interface {
	transaction(f func(tx *gorm.DB) error) error
	findByUserName(ctx context.Context, username string) (*model.User, error)
	findByEmail(ctx context.Context, email string) (*model.User, error)
	saveUser(ctx context.Context, tx *gorm.DB, u *model.User) error
}
