package repo

import (
	"netease-kit/nemo/internal/db"
	"netease-kit/nemo/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{DB: db.DB}
}

func (r *UserRepository) SelectByUserUuid(userUuid string) (*model.User, error) {
	var user model.User
	result := r.DB.Where("user_uuid = ?", userUuid).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) Insert(user *model.User) error {
	return r.DB.Create(user).Error
}
