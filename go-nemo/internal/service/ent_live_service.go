package service

import (
	"errors"
	"netease-kit/nemo/internal/client"
	"netease-kit/nemo/internal/model"
	"netease-kit/nemo/internal/repo"

	"github.com/google/uuid"
)

type EntLiveService struct {
	LiveRecordRepo *repo.LiveRecordRepository
	LiveRewardRepo *repo.LiveRewardRepository
	NeRoomClient   *client.NeRoomClient
	NimClient      *client.NimClient
}

func NewEntLiveService() *EntLiveService {
	return &EntLiveService{
		LiveRecordRepo: repo.NewLiveRecordRepository(),
		LiveRewardRepo: repo.NewLiveRewardRepository(),
		NeRoomClient:   client.NewNeRoomClient(),
		NimClient:      client.NewNimClient(),
	}
}

func (s *EntLiveService) CreateLive(userUuid string, liveTopic string, cover string, liveType int) (*model.LiveRecord, error) {
	// 1. Check if user already has a live room
	existRecord, err := s.LiveRecordRepo.SelectByUserUuidAndType(userUuid, liveType)
	if err != nil {
		return nil, err
	}
	if existRecord != nil {
		// If exists, return it (or close it and create new one, but Java version returns it)
		return existRecord, nil
	}

	// 2. Generate Room UUID
	roomUuid := uuid.New().String()

	// 3. Create NeRoom
	// Assuming default template ID for now, or map based on liveType
	templateId := int64(1001) // Example template ID
	param := client.CreateNeRoomParam{
		RoomName:   liveTopic,
		RoomUuid:   roomUuid,
		TemplateId: templateId,
		RoomSeatConfig: &client.RoomSeatConfig{
			SeatCount:  8,
			ApplyMode:  1,
			InviteMode: 1,
		},
		RoomConfig: &client.RoomConfig{
			Resource: &client.ResourceConfig{
				Rtc:      true,
				Chatroom: true,
				Live:     true,
			},
		},
	}

	neRoomDto, err := s.NeRoomClient.CreateNeRoom(param)
	if err != nil {
		return nil, err
	}

	// 4. Save LiveRecord
	record := &model.LiveRecord{
		UserUuid:      userUuid,
		LiveTopic:     liveTopic,
		Cover:         cover,
		LiveType:      liveType,
		RoomUuid:      roomUuid,
		RoomArchiveId: neRoomDto.RoomArchiveId,
		Status:        1,
		Live:          1,
	}

	err = s.LiveRecordRepo.Insert(record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *EntLiveService) CloseLive(userUuid string, roomArchiveId string) error {
	record, err := s.LiveRecordRepo.SelectByRoomArchiveId(roomArchiveId)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.New("live record not found")
	}

	if record.UserUuid != userUuid {
		return errors.New("permission denied")
	}

	// Call NeRoom to destroy room
	err = s.NeRoomClient.DeleteNeRoom(roomArchiveId)
	if err != nil {
		// Log error but continue to update DB? Or fail?
		// Java version throws exception.
		return err
	}

	// Update DB
	err = s.LiveRecordRepo.UpdateStatus(record.Id, -1, -1)
	if err != nil {
		return err
	}

	return nil
}

func (s *EntLiveService) GetLiveList(liveType int, excludeUserUuid string) ([]*model.LiveRecord, error) {
	return s.LiveRecordRepo.GetLivingRecords(liveType, excludeUserUuid)
}

func (s *EntLiveService) GetLiveInfo(liveRecordId uint64) (*model.LiveRecord, error) {
	return s.LiveRecordRepo.SelectById(liveRecordId)
}
