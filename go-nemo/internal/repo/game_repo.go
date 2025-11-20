package repo

import (
	"netease-kit/nemo/internal/db"
	"netease-kit/nemo/internal/model"

	"gorm.io/gorm"
)

type GameRecordRepository struct {
	DB *gorm.DB
}

func NewGameRecordRepository() *GameRecordRepository {
	return &GameRecordRepository{DB: db.DB}
}

func (r *GameRecordRepository) Insert(record *model.GameRecord) error {
	return r.DB.Create(record).Error
}

func (r *GameRecordRepository) SelectByRoomUuid(roomUuid string) (*model.GameRecord, error) {
	var record model.GameRecord
	result := r.DB.Where("room_uuid = ? AND game_status != 2", roomUuid).First(&record)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &record, nil
}

func (r *GameRecordRepository) UpdateStatus(id uint64, status int) error {
	return r.DB.Model(&model.GameRecord{}).Where("id = ?", id).Update("game_status", status).Error
}

type GameMemberRepository struct {
	DB *gorm.DB
}

func NewGameMemberRepository() *GameMemberRepository {
	return &GameMemberRepository{DB: db.DB}
}

func (r *GameMemberRepository) Insert(member *model.GameMember) error {
	return r.DB.Create(member).Error
}

func (r *GameMemberRepository) SelectByGameRecordId(gameRecordId uint64) ([]model.GameMember, error) {
	var members []model.GameMember
	result := r.DB.Where("game_record_id = ?", gameRecordId).Find(&members)
	return members, result.Error
}

func (r *GameMemberRepository) SelectByUserUuidAndGameRecordId(userUuid string, gameRecordId uint64) (*model.GameMember, error) {
	var member model.GameMember
	result := r.DB.Where("user_uuid = ? AND game_record_id = ?", userUuid, gameRecordId).First(&member)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &member, nil
}

func (r *GameMemberRepository) UpdateStatus(id uint64, status int) error {
	return r.DB.Model(&model.GameMember{}).Where("id = ?", id).Update("status", status).Error
}
