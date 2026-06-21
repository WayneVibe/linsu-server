package auths

import (
	"context"
	"model"

	"gorm.io/gorm"

	"github.com/setcreed/hade-kit/gorms"
)

type models struct {
	db *gorm.DB
}

func (m *models) saveUser(ctx context.Context, tx *gorm.DB, u *model.User) error {
	if tx == nil {
		tx = m.db
	}
	return tx.WithContext(ctx).Create(u).Error
}

func (m *models) transaction(f func(tx *gorm.DB) error) error {
	return m.db.Transaction(f)
}

func (m *models) findByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := m.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &user, err
}

func (m *models) findByUserName(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := m.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &user, err
}

func newModel(db *gorm.DB) *models {
	return &models{db: db}
}
