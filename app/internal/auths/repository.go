package auths

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"model"
)

type repository interface {
	transaction(ctx context.Context, f func(tx *gorm.DB) error) error
	findByUserName(ctx context.Context, username string) (*model.User, error)
	findByEmail(ctx context.Context, email string) (*model.User, error)
	saveUser(ctx context.Context, tx *gorm.DB, u *model.User) error
	findById(ctx context.Context, id uuid.UUID) (*model.User, error)
	updateUser(ctx context.Context, tx *gorm.DB, u *model.User) error
}
