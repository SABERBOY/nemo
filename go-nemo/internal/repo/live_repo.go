package repo

import (
	"netease-kit/nemo/internal/db"
	"netease-kit/nemo/internal/model"

	"gorm.io/gorm"
)

type LiveRecordRepository struct {
	DB *gorm.DB
}

func NewLiveRecordRepository() *LiveRecordRepository {
	return &LiveRecordRepository{DB: db.DB}
}

func (r *LiveRecordRepository) Insert(record *model.LiveRecord) error {
	return r.DB.Create(record).Error
}

func (r *LiveRecordRepository) SelectByRoomArchiveId(roomArchiveId string) (*model.LiveRecord, error) {
	var record model.LiveRecord
	result := r.DB.Where("room_archive_id = ?", roomArchiveId).First(&record)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &record, nil
}

func (r *LiveRecordRepository) UpdateStatus(id uint64, status int, live int) error {
	return r.DB.Model(&model.LiveRecord{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
		"live":   live,
	}).Error
}

func (r *LiveRecordRepository) SelectByUserUuidAndType(userUuid string, liveType int) (*model.LiveRecord, error) {
	var record model.LiveRecord
	result := r.DB.Where("user_uuid = ? AND live_type = ? AND status = 1", userUuid, liveType).First(&record)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &record, nil
}

func (r *LiveRecordRepository) GetLivingRecords(liveType int, excludeUserUuid string) ([]*model.LiveRecord, error) {
	var records []*model.LiveRecord
	query := r.DB.Where("live = 1 AND status = 1")
	if liveType != 0 {
		query = query.Where("live_type = ?", liveType)
	}
	if excludeUserUuid != "" {
		query = query.Where("user_uuid != ?", excludeUserUuid)
	}
	err := query.Order("id desc").Find(&records).Error
	return records, err
}

func (r *LiveRecordRepository) SelectById(id uint64) (*model.LiveRecord, error) {
	var record model.LiveRecord
	result := r.DB.First(&record, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &record, nil
}

type LiveRewardRepository struct {
	DB *gorm.DB
}

func NewLiveRewardRepository() *LiveRewardRepository {
	return &LiveRewardRepository{DB: db.DB}
}

func (r *LiveRewardRepository) Insert(reward *model.LiveReward) error {
	return r.DB.Create(reward).Error
}
