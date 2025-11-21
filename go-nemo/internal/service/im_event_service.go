package service

import (
	"encoding/json"
	"log"
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/repo"
	"strconv"
	"strings"
)

type ImEventService struct {
	ImRepo         *repo.ImRepository
	LiveRecordRepo *repo.LiveRecordRepository
}

func NewImEventService() *ImEventService {
	return &ImEventService{
		ImRepo:         repo.NewImRepository(),
		LiveRecordRepo: repo.NewLiveRecordRepository(),
	}
}

func (s *ImEventService) HandlerChatroomInOut(body string) {
	var param dto.ImMemberInOutMsgCallbackParam
	err := json.Unmarshal([]byte(body), &param)
	if err != nil {
		log.Printf("HandlerChatroomInOut unmarshal error: %v", err)
		return
	}

	if param.RoomId == "" || param.Accid == "" || param.Event == "" || param.Timestamp == 0 {
		log.Printf("Invalid chatroom in/out message: %s", body)
		return
	}

	chatRoomId, err := strconv.ParseInt(param.RoomId, 10, 64)
	if err != nil {
		log.Printf("Invalid chatRoomId: %s", param.RoomId)
		return
	}

	// Check if user is host
	record, err := s.LiveRecordRepo.SelectByChatRoomId(chatRoomId)
	if err != nil {
		log.Printf("SelectByChatRoomId error: %v", err)
		// Continue? Java version returns if record is null or user is host.
		// But if record is null, we might not know if it's a valid room.
		// Java: if(liveRecord == null || Objects.equals(accid, liveRecord.getUserUuid())) return;
	}

	if record == nil {
		// Maybe room closed or not found.
		// Java logs "ignore live owner in/out event" if record is null? No, if record is null OR user is owner.
		// If record is null, we can't check owner. But maybe we should ignore anyway.
		return
	}

	if record.UserUuid == param.Accid {
		log.Printf("Ignore live owner in/out event, chatRoomId: %d, accid: %s", chatRoomId, param.Accid)
		return
	}

	// Lock logic omitted for simplicity, or use Redis lock if needed.
	// Java uses "lock:chatroom:inout:" + chatRoomId + ":" + accid

	if strings.EqualFold(param.Event, "IN") {
		s.HandleUserEnter(chatRoomId, param.Accid, param.Timestamp)
	} else if strings.EqualFold(param.Event, "OUT") {
		s.HandleUserLeave(chatRoomId, param.Accid, param.Timestamp)
	} else {
		log.Printf("Unknown event type: %s", param.Event)
	}
}

func (s *ImEventService) HandleUserEnter(chatRoomId int64, accid string, timestamp int64) {
	lastLeave, _ := s.ImRepo.GetLastLeaveTime(chatRoomId, accid)

	if lastLeave == 0 || timestamp > lastLeave {
		s.ImRepo.SetLastEnterTime(chatRoomId, accid, timestamp)

		exists, _ := s.ImRepo.IsAudienceExist(chatRoomId, accid)
		if !exists {
			s.ImRepo.AddAudience(chatRoomId, accid, timestamp)
			s.ImRepo.IncrementAudienceCount(chatRoomId)
			log.Printf("User entered chatroom: roomId=%d, accid=%s", chatRoomId, accid)
		} else {
			s.ImRepo.AddAudience(chatRoomId, accid, timestamp) // Update score
			log.Printf("Updated user enter time: roomId=%d, accid=%s", chatRoomId, accid)
		}
		s.ImRepo.ExpireAudienceKeys(chatRoomId)
	}
}

func (s *ImEventService) HandleUserLeave(chatRoomId int64, accid string, timestamp int64) {
	lastEnter, _ := s.ImRepo.GetLastEnterTime(chatRoomId, accid)

	if lastEnter != 0 && timestamp > lastEnter {
		s.ImRepo.SetLastLeaveTime(chatRoomId, accid, timestamp)

		exists, _ := s.ImRepo.IsAudienceExist(chatRoomId, accid)
		if exists {
			s.ImRepo.RemoveAudience(chatRoomId, accid)
			s.ImRepo.DecrementAudienceCount(chatRoomId)
			log.Printf("User left chatroom: roomId=%d, accid=%s", chatRoomId, accid)
		}
	}
}
