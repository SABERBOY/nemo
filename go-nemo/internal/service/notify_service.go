package service

import (
	"encoding/json"
	"log"
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/repo"
	"time"
)

type NotifyService struct {
	LiveRecordRepo *repo.LiveRecordRepository
	NeRoomRepo     *repo.NeRoomRepository
	EntLiveService *EntLiveService
	MusicService   *MusicPlayService
	KtvService     *KtvService
}

func NewNotifyService(
	entLiveService *EntLiveService,
	musicService *MusicPlayService,
	ktvService *KtvService,
) *NotifyService {
	return &NotifyService{
		LiveRecordRepo: repo.NewLiveRecordRepository(),
		NeRoomRepo:     repo.NewNeRoomRepository(),
		EntLiveService: entLiveService,
		MusicService:   musicService,
		KtvService:     ktvService,
	}
}

func (s *NotifyService) HandleNotify(body string) {
	var notifyBody dto.NeRoomNotifyBody
	err := json.Unmarshal([]byte(body), &notifyBody)
	if err != nil {
		log.Printf("HandleNotify unmarshal error: %v", err)
		return
	}

	dataBytes, _ := json.Marshal(notifyBody.Body)

	switch notifyBody.EventType {
	case "CREATE_ROOM":
		// TODO
	case "CLOSE_ROOM":
		var param dto.CloseRoomEventNotify
		json.Unmarshal(dataBytes, &param)
		s.HandlerCloseRoom(param)
	case "USER_JOIN_ROOM":
		var param dto.JoinRoomEventNotify
		json.Unmarshal(dataBytes, &param)
		s.HandlerUserJoinRoom(param)
	case "USER_LEAVE_ROOM":
		var param dto.LeaveRoomEventNotify
		json.Unmarshal(dataBytes, &param)
		s.HandlerUserLeaveRoom(param)
	case "USER_ON_SEAT":
		var param dto.UserOnSeatNotifyDto
		json.Unmarshal(dataBytes, &param)
		s.HandlerUserOnSeat(param)
	case "USER_OFF_SEAT":
		var param dto.UserOffSeatNotifyDto
		json.Unmarshal(dataBytes, &param)
		s.HandlerUserOffSeat(param)
	}
}

func (s *NotifyService) HandlerCloseRoom(param dto.CloseRoomEventNotify) {
	record, err := s.LiveRecordRepo.SelectByRoomArchiveId(param.RoomArchiveId)
	if err != nil || record == nil {
		return
	}
	if record.Live == -1 { // LIVE_CLOSE
		return
	}

	// Update live state to closed
	s.LiveRecordRepo.UpdateStatus(record.Id, -1, -1)

	// Clean music info
	s.MusicService.CleanPlayerMusicInfo(record.Id)

	// Clean order songs
	s.KtvService.CleanOrderSongs(record.Id)

	// Clean sing info (if KTV)
	if record.LiveType == 2 { // KTV
		s.KtvService.CleanSingInfo(record.RoomUuid)
	}

	// Game room close event (omitted for now)
}

func (s *NotifyService) HandlerUserJoinRoom(param dto.JoinRoomEventNotify) {
	record, err := s.LiveRecordRepo.SelectByRoomArchiveId(param.RoomArchiveId)
	if err != nil || record == nil {
		return
	}

	for _, user := range param.Users {
		s.NeRoomRepo.AddMember(param.RoomArchiveId, user)

		// If user is host and live is not started or paused, set to LIVE
		if user.UserUuid == record.UserUuid && (record.Live == 0 || record.Live == 2) {
			if record.LiveType != 3 { // Not PK Live
				s.LiveRecordRepo.UpdateStatus(record.Id, 1, 1) // LIVE
			} else {
				// Update PK Live Layout (omitted)
			}
		}
	}
	s.NeRoomRepo.ExpireMemberTable(param.RoomArchiveId, 7*24*time.Hour)
}

func (s *NotifyService) HandlerUserLeaveRoom(param dto.LeaveRoomEventNotify) {
	record, err := s.LiveRecordRepo.SelectByRoomArchiveId(param.RoomArchiveId)
	if err != nil || record == nil {
		return
	}

	for _, user := range param.Users {
		s.NeRoomRepo.RemoveMember(param.RoomArchiveId, user.UserUuid)

		if user.UserUuid == record.UserUuid {
			// If host leaves and not PK live, close room
			if record.LiveType != 3 {
				s.EntLiveService.CloseLive(user.UserUuid, param.RoomArchiveId)
			}
		}

		// KTV: Clean user order songs
		if record.LiveType == 2 {
			oldOrderSongs, _ := s.KtvService.GetOrderSongs(record.Id)
			s.KtvService.CleanUserOrderSongs(record.Id, user.UserUuid)
			s.KtvService.EndSingWhenMemberOut(record.RoomUuid, user.UserUuid, oldOrderSongs)
		}
	}
}

func (s *NotifyService) HandlerUserOnSeat(param dto.UserOnSeatNotifyDto) {
	record, err := s.LiveRecordRepo.SelectByRoomArchiveId(param.RoomArchiveId)
	if err != nil || record == nil {
		return
	}

	s.NeRoomRepo.AddSeatUser(param.RoomArchiveId, param.Index, param.User)

	// Update PK Live Layout (omitted)
}

func (s *NotifyService) HandlerUserOffSeat(param dto.UserOffSeatNotifyDto) {
	s.NeRoomRepo.RemoveSeatUser(param.RoomArchiveId, param.Index)

	record, err := s.LiveRecordRepo.SelectByRoomArchiveId(param.RoomArchiveId)
	if err != nil || record == nil {
		return
	}

	// Update PK Live Layout (omitted)

	if record.LiveType == 2 { // KTV
		oldOrderSongs, _ := s.KtvService.GetOrderSongs(record.Id)
		s.KtvService.CleanUserOrderSongs(record.Id, param.UserUuid)
		s.KtvService.EndSingWhenMemberOut(record.RoomUuid, param.UserUuid, oldOrderSongs)
	}
}
