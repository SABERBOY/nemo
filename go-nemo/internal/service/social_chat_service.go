package service

import (
	"context"
	"encoding/json"
	"fmt"
	"netease-kit/nemo/internal/client"
	"netease-kit/nemo/internal/config"
	"netease-kit/nemo/internal/db"
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/model"
	"netease-kit/nemo/internal/repo"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	RedisKeyPrefixOneOneOnlineUser      = "nemo:one:one:online:user"
	RedisKeyPrefixOneOneOnlineUserOrder = "nemo:one:one:online:user:order"
	RedisKeyPrefixRtcRecord             = "nemo:rtc:record"
	RedisKeyPrefixOneOneChatRtcUser     = "nemo:one:one:chat:rtc:user:record"
)

type SocialChatService struct {
	UserRepo       *repo.UserRepository
	GiftRepo       *repo.GiftRepository
	UserRewardRepo *repo.UserRewardRepository
	NimClient      *client.NimClient
	RDB            *redis.Client
}

func NewSocialChatService() *SocialChatService {
	return &SocialChatService{
		UserRepo:       repo.NewUserRepository(),
		GiftRepo:       repo.NewGiftRepository(),
		UserRewardRepo: repo.NewUserRewardRepository(),
		NimClient:      client.NewNimClient(),
		RDB:            db.RDB,
	}
}

func (s *SocialChatService) Reporter(appKey, userUuid, deviceId string) error {
	ctx := context.Background()
	userKey := fmt.Sprintf("%s:%s:%s", RedisKeyPrefixOneOneOnlineUser, appKey, userUuid)
	orderKey := fmt.Sprintf("%s:%s", RedisKeyPrefixOneOneOnlineUserOrder, appKey)
	now := time.Now().UnixMilli()

	// Check if user exists in Redis
	val, err := s.RDB.Get(ctx, userKey).Result()
	var onLineUser dto.OnLineUserDto

	if err == redis.Nil {
		// Not found, get from DB
		user, err := s.UserRepo.SelectByUserUuid(userUuid)
		if err != nil {
			return err
		}
		if user == nil {
			return fmt.Errorf("user not found: %s", userUuid)
		}
		onLineUser = dto.OnLineUserDto{
			UserUuid:        user.UserUuid,
			UserName:        user.UserName,
			Icon:            user.Icon,
			Mobile:          user.Mobile,
			FirstReportTime: now,
		}
	} else if err != nil {
		return err
	} else {
		json.Unmarshal([]byte(val), &onLineUser)
	}

	onLineUser.LastReportTime = now
	jsonBytes, _ := json.Marshal(onLineUser)

	pipe := s.RDB.Pipeline()
	pipe.Set(ctx, userKey, string(jsonBytes), 60*time.Second)
	pipe.ZAdd(ctx, orderKey, &redis.Z{Score: float64(now), Member: userUuid})
	_, err = pipe.Exec(ctx)
	return err
}

func (s *SocialChatService) GetOnLineUser(appKey string, pageNum, pageSize int) ([]dto.OnLineUserDto, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", RedisKeyPrefixOneOneOnlineUser, appKey)
	orderKey := fmt.Sprintf("%s:%s", RedisKeyPrefixOneOneOnlineUserOrder, appKey)

	start := int64((pageNum - 1) * pageSize)
	end := int64(pageNum*pageSize - 1)

	userUuids, err := s.RDB.ZRange(ctx, orderKey, start, end).Result()
	if err != nil {
		return nil, err
	}

	if len(userUuids) == 0 {
		return []dto.OnLineUserDto{}, nil
	}

	var keys []string
	for _, uuid := range userUuids {
		keys = append(keys, fmt.Sprintf("%s:%s", key, uuid))
	}

	vals, err := s.RDB.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	var users []dto.OnLineUserDto
	for _, val := range vals {
		if val == nil {
			continue
		}
		var user dto.OnLineUserDto
		json.Unmarshal([]byte(val.(string)), &user)
		users = append(users, user)
	}

	return users, nil
}

func (s *SocialChatService) UserReward(userRewardDto dto.UserRewardDto) error {
	gift, err := s.GiftRepo.SelectByPrimaryKey(userRewardDto.GiftId)
	if err != nil {
		return err
	}
	if gift == nil {
		return fmt.Errorf("gift not found")
	}

	user, err := s.UserRepo.SelectByUserUuid(userRewardDto.UserUuid)
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}

	targetUser, err := s.UserRepo.SelectByUserUuid(userRewardDto.Target)
	if err != nil || targetUser == nil {
		return fmt.Errorf("target user not found")
	}

	reward := &model.UserReward{
		UserUuid:  userRewardDto.UserUuid,
		GiftId:    userRewardDto.GiftId,
		GiftCount: userRewardDto.GiftCount,
		CloudCoin: gift.CloudCoin,
		Target:    userRewardDto.Target,
	}

	err = s.UserRewardRepo.Insert(reward)
	if err != nil {
		return err
	}

	// Send Notification
	s.notifyRewardMessage(reward, targetUser, user)

	return nil
}

func (s *SocialChatService) notifyRewardMessage(reward *model.UserReward, targetUser, user *model.User) {
	// 1. Send to Target
	msg := dto.RewardMessage{
		SenderUserUuid: user.UserUuid,
		TargetUserUuid: targetUser.UserUuid,
		GiftId:         reward.GiftId,
		GiftCount:      reward.GiftCount,
		CloudCoin:      reward.CloudCoin,
	}
	event := dto.EventDto{
		Data: msg,
		Type: 10, // EventTypeEnum.SOCIAL_CHAT_USER_REWARD
	}
	s.NimClient.SendImCustomMsg(user.UserUuid, targetUser.UserUuid, 0, event)

	// 2. Send from Assistant
	assistAccid := config.AppConfig.Business.YunxinAssistAccid
	// Simplified message for assistant
	s.NimClient.SendImCustomMsg(assistAccid, targetUser.UserUuid, 0, event)
}

func (s *SocialChatService) SaveRtcRecord(appKey string, param dto.RtcRoomNotifyParam) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", RedisKeyPrefixRtcRecord, appKey)

	info := dto.RtcRoomInfoDto{
		ChannelId:   param.ChannelId,
		ChannelName: param.ChannelName,
		Status:      param.Status,
		AppKey:      appKey,
	}

	jsonBytes, _ := json.Marshal(info)
	return s.RDB.HSet(ctx, key, fmt.Sprintf("%d", param.ChannelId), string(jsonBytes)).Err()
}
