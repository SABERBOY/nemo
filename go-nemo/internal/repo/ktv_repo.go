package repo

import (
	"netease-kit/nemo/internal/db"
	"netease-kit/nemo/internal/model"
	"time"

	"gorm.io/gorm"
)

type OrderSongRepository struct {
	DB *gorm.DB
}

func NewOrderSongRepository() *OrderSongRepository {
	return &OrderSongRepository{DB: db.DB}
}

func (r *OrderSongRepository) Insert(orderSong *model.OrderSong) error {
	return r.DB.Create(orderSong).Error
}

func (r *OrderSongRepository) SelectByLiveRecordId(liveRecordId uint64) ([]*model.OrderSong, error) {
	var songs []*model.OrderSong
	err := r.DB.Where("live_record_id = ? AND status >= 0", liveRecordId).Order("set_top_time desc, id asc").Find(&songs).Error
	return songs, err
}

func (r *OrderSongRepository) SelectById(id uint64) (*model.OrderSong, error) {
	var song model.OrderSong
	err := r.DB.First(&song, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &song, nil
}

func (r *OrderSongRepository) UpdateStatus(id uint64, status int) error {
	return r.DB.Model(&model.OrderSong{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderSongRepository) SetTop(id uint64) error {
	return r.DB.Model(&model.OrderSong{}).Where("id = ?", id).Update("set_top_time", time.Now().UnixMilli()).Error
}

func (r *OrderSongRepository) Delete(id uint64) error {
	return r.DB.Model(&model.OrderSong{}).Where("id = ?", id).Update("status", -1).Error
}

func (r *OrderSongRepository) GetNextSong(liveRecordId uint64) (*model.OrderSong, error) {
	var song model.OrderSong
	err := r.DB.Where("live_record_id = ? AND status = 0", liveRecordId).Order("set_top_time desc, id asc").First(&song).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &song, nil
}

type ChorusRecordRepository struct {
	DB *gorm.DB
}

func NewChorusRecordRepository() *ChorusRecordRepository {
	return &ChorusRecordRepository{DB: db.DB}
}

func (r *ChorusRecordRepository) Insert(record *model.ChorusRecord) error {
	return r.DB.Create(record).Error
}

func (r *ChorusRecordRepository) SelectByChorusId(chorusId string) (*model.ChorusRecord, error) {
	var record model.ChorusRecord
	err := r.DB.Where("chorus_id = ?", chorusId).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *ChorusRecordRepository) UpdateState(id uint64, state int) error {
	return r.DB.Model(&model.ChorusRecord{}).Where("id = ?", id).Update("state", state).Error
}
