package repo

import (
	"netease-kit/nemo/internal/db"
	"netease-kit/nemo/internal/model"

	"gorm.io/gorm"
)

type GiftRepository struct {
	DB *gorm.DB
}

func NewGiftRepository() *GiftRepository {
	return &GiftRepository{DB: db.DB}
}

func (r *GiftRepository) SelectByPrimaryKey(id uint64) (*model.Gift, error) {
	var gift model.Gift
	result := r.DB.First(&gift, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &gift, nil
}

type UserRewardRepository struct {
	DB *gorm.DB
}

func NewUserRewardRepository() *UserRewardRepository {
	return &UserRewardRepository{DB: db.DB}
}

func (r *UserRewardRepository) Insert(reward *model.UserReward) error {
	return r.DB.Create(reward).Error
}
