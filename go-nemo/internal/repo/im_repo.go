package repo

import (
	"context"
	"fmt"
	"netease-kit/nemo/internal/db"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type ImRepository struct{}

func NewImRepository() *ImRepository {
	return &ImRepository{}
}

const (
	KeyAudienceList      = "chatroom:audience:list:"
	KeyAudienceCount     = "chatroom:audience:count:"
	KeyAudienceLastEnter = "chatroom:audience:last_enter:"
	KeyAudienceLastLeave = "chatroom:audience:last_leave:"
	ExpireTime           = 24 * time.Hour
)

func (r *ImRepository) GetLastEnterTime(chatRoomId int64, accid string) (int64, error) {
	key := fmt.Sprintf("%s%d:%s", KeyAudienceLastEnter, chatRoomId, accid)
	val, err := db.RDB.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

func (r *ImRepository) SetLastEnterTime(chatRoomId int64, accid string, timestamp int64) error {
	key := fmt.Sprintf("%s%d:%s", KeyAudienceLastEnter, chatRoomId, accid)
	return db.RDB.Set(context.Background(), key, timestamp, ExpireTime).Err()
}

func (r *ImRepository) GetLastLeaveTime(chatRoomId int64, accid string) (int64, error) {
	key := fmt.Sprintf("%s%d:%s", KeyAudienceLastLeave, chatRoomId, accid)
	val, err := db.RDB.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

func (r *ImRepository) SetLastLeaveTime(chatRoomId int64, accid string, timestamp int64) error {
	key := fmt.Sprintf("%s%d:%s", KeyAudienceLastLeave, chatRoomId, accid)
	return db.RDB.Set(context.Background(), key, timestamp, ExpireTime).Err()
}

func (r *ImRepository) AddAudience(chatRoomId int64, accid string, timestamp int64) error {
	key := fmt.Sprintf("%s%d", KeyAudienceList, chatRoomId)
	return db.RDB.ZAdd(context.Background(), key, &redis.Z{
		Score:  float64(timestamp),
		Member: accid,
	}).Err()
}

func (r *ImRepository) RemoveAudience(chatRoomId int64, accid string) error {
	key := fmt.Sprintf("%s%d", KeyAudienceList, chatRoomId)
	return db.RDB.ZRem(context.Background(), key, accid).Err()
}

func (r *ImRepository) IsAudienceExist(chatRoomId int64, accid string) (bool, error) {
	key := fmt.Sprintf("%s%d", KeyAudienceList, chatRoomId)
	_, err := db.RDB.ZScore(context.Background(), key, accid).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ImRepository) IncrementAudienceCount(chatRoomId int64) error {
	key := fmt.Sprintf("%s%d", KeyAudienceCount, chatRoomId)
	return db.RDB.Incr(context.Background(), key).Err()
}

func (r *ImRepository) DecrementAudienceCount(chatRoomId int64) error {
	key := fmt.Sprintf("%s%d", KeyAudienceCount, chatRoomId)
	val, err := db.RDB.Decr(context.Background(), key).Result()
	if err != nil {
		return err
	}
	if val < 0 {
		db.RDB.Set(context.Background(), key, 0, 0)
	}
	return nil
}

func (r *ImRepository) ExpireAudienceKeys(chatRoomId int64) {
	listKey := fmt.Sprintf("%s%d", KeyAudienceList, chatRoomId)
	countKey := fmt.Sprintf("%s%d", KeyAudienceCount, chatRoomId)
	db.RDB.Expire(context.Background(), listKey, ExpireTime)
	db.RDB.Expire(context.Background(), countKey, ExpireTime)
}
