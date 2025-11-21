package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"netease-kit/nemo/internal/db"
	"netease-kit/nemo/internal/dto"
	"time"

	"github.com/go-redis/redis/v8"
)

type MusicPlayerRepository struct {
	RDB *redis.Client
}

func NewMusicPlayerRepository() *MusicPlayerRepository {
	return &MusicPlayerRepository{RDB: db.RDB}
}

func (r *MusicPlayerRepository) GetMusicPlayerInfo(liveRecordId uint64) (*dto.PlayDetailInfoDto, error) {
	key := fmt.Sprintf("nemo:ent:music:player:%d", liveRecordId)
	val, err := r.RDB.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var info dto.PlayDetailInfoDto
	err = json.Unmarshal([]byte(val), &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *MusicPlayerRepository) PutMusicPlayerInfo(liveRecordId uint64, info *dto.PlayDetailInfoDto) error {
	key := fmt.Sprintf("nemo:ent:music:player:%d", liveRecordId)
	val, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return r.RDB.Set(context.Background(), key, val, 24*time.Hour).Err()
}

func (r *MusicPlayerRepository) DeleteMusicPlayerInfo(liveRecordId uint64) error {
	key := fmt.Sprintf("nemo:ent:music:player:%d", liveRecordId)
	return r.RDB.Del(context.Background(), key).Err()
}

func (r *MusicPlayerRepository) MusicReady(liveRecordId uint64, orderId uint64, userUuid string) error {
	key := fmt.Sprintf("nemo:ent:music:ready:%d:%d", liveRecordId, orderId)
	return r.RDB.SAdd(context.Background(), key, userUuid).Err()
}

func (r *MusicPlayerRepository) IsMusicReady(liveRecordId uint64, orderId uint64, userUuid string) (bool, error) {
	key := fmt.Sprintf("nemo:ent:music:ready:%d:%d", liveRecordId, orderId)
	return r.RDB.SIsMember(context.Background(), key, userUuid).Result()
}

func (r *MusicPlayerRepository) DeleteMusicReady(liveRecordId uint64, orderId uint64) error {
	key := fmt.Sprintf("nemo:ent:music:ready:%d:%d", liveRecordId, orderId)
	return r.RDB.Del(context.Background(), key).Err()
}
