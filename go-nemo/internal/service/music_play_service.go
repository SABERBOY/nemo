package service

import (
	"encoding/json"
	"errors"
	"netease-kit/nemo/internal/client"
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/model"
	"netease-kit/nemo/internal/repo"
)

type MusicPlayService struct {
	MusicPlayerRepo *repo.MusicPlayerRepository
	LiveRecordRepo  *repo.LiveRecordRepository
	OrderSongRepo   *repo.OrderSongRepository
	UserRepo        *repo.UserRepository
	NeRoomClient    *client.NeRoomClient
}

func NewMusicPlayService() *MusicPlayService {
	return &MusicPlayService{
		MusicPlayerRepo: repo.NewMusicPlayerRepository(),
		LiveRecordRepo:  repo.NewLiveRecordRepository(),
		OrderSongRepo:   repo.NewOrderSongRepository(),
		UserRepo:        repo.NewUserRepository(),
		NeRoomClient:    client.NewNeRoomClient(),
	}
}

func (s *MusicPlayService) GetPlayMusicInfo(liveRecordId uint64) (*dto.PlayDetailInfoDto, error) {
	liveRecord, err := s.LiveRecordRepo.SelectById(liveRecordId)
	if err != nil || liveRecord == nil {
		return nil, errors.New("live record not found")
	}
	if liveRecord.Live != 1 {
		return nil, errors.New("anchor not living")
	}

	return s.MusicPlayerRepo.GetMusicPlayerInfo(liveRecordId)
}

func (s *MusicPlayService) MusicReady(liveRecordId uint64, orderId uint64, userUuid string) error {
	liveRecord, err := s.LiveRecordRepo.SelectById(liveRecordId)
	if err != nil || liveRecord == nil {
		return errors.New("live record not found")
	}
	if liveRecord.Live != 1 {
		return errors.New("anchor not living")
	}

	// Check if user in room (Skip for now or implement NeRoom check)

	orderSong, err := s.OrderSongRepo.SelectById(orderId)
	if err != nil || orderSong == nil {
		return errors.New("order song not found")
	}
	if orderSong.LiveRecordId != liveRecordId {
		return errors.New("order song mismatch")
	}

	// Simplified logic: Direct play for now (Chat Room logic)
	// TODO: Implement Listen Together logic with delay queue
	return s.musicPlay(liveRecordId, userUuid, orderSong)
}

func (s *MusicPlayService) MusicAction(userUuid string, param dto.MusicActionParam) error {
	liveRecordId := param.LiveRecordId
	liveRecord, err := s.LiveRecordRepo.SelectById(liveRecordId)
	if err != nil || liveRecord == nil {
		return errors.New("live record not found")
	}
	if liveRecord.Live != 1 {
		return errors.New("anchor not living")
	}

	info, err := s.MusicPlayerRepo.GetMusicPlayerInfo(liveRecordId)
	if err != nil || info == nil {
		return errors.New("music info not found")
	}

	// Update status
	status := 0
	switch param.Action {
	case 1: // Play
		status = 1
	case 2: // Pause
		status = 2
	case 3: // Resume
		status = 1
	case 4: // End
		// Handle end
	}
	info.MusicStatus = status
	s.MusicPlayerRepo.PutMusicPlayerInfo(liveRecordId, info)

	if param.Action == 1 { // Play
		s.setSongPlaying(info.OrderId)
		s.resetOtherOrderSongStatus(liveRecordId, info.OrderId)
	}

	// Send Notify
	operator := s.getOperatorInfoDto(userUuid)
	event := dto.MusicPlayNotifyEventDto{
		Operator:      operator,
		PlayMusicInfo: info,
	}

	// EventType: 1: Play, 2: Pause, 3: Resume
	eventType := 0
	if param.Action == 1 {
		eventType = 10 // ENT_MUSIC_PLAY (Example code)
	} else if param.Action == 2 {
		eventType = 11 // ENT_MUSIC_PAUSE
	} else if param.Action == 3 {
		eventType = 12 // ENT_MUSIC_RESUME
	}

	s.sendMusicPlayerMessage(event, info.RoomUuid, eventType)
	return nil
}

func (s *MusicPlayService) musicPlay(liveRecordId uint64, userUuid string, orderSong *model.OrderSong) error {
	operator := s.getOperatorInfoDto(userUuid)

	info := &dto.PlayDetailInfoDto{
		OrderId:     orderSong.Id,
		SongId:      orderSong.SongId,
		SongName:    orderSong.SongName,
		SongCover:   orderSong.SongCover,
		Singer:      orderSong.Singer,
		SingerCover: orderSong.SingerCover,
		SongTime:    orderSong.SongTime,
		Channel:     orderSong.Channel,
		MusicStatus: 1, // Play
		RoomUuid:    orderSong.RoomUuid,
		UserUuid:    orderSong.UserUuid,
	}
	// Fill user info for song owner
	songOwner, _ := s.UserRepo.SelectByUserUuid(orderSong.UserUuid)
	if songOwner != nil {
		info.UserName = songOwner.UserName
		info.UserIcon = songOwner.Icon
	}

	s.MusicPlayerRepo.PutMusicPlayerInfo(liveRecordId, info)
	s.setSongPlaying(orderSong.Id)
	s.resetOtherOrderSongStatus(liveRecordId, orderSong.Id)

	event := dto.MusicPlayNotifyEventDto{
		Operator:      operator,
		PlayMusicInfo: info,
	}
	s.sendMusicPlayerMessage(event, info.RoomUuid, 10) // ENT_MUSIC_PLAY
	return nil
}

func (s *MusicPlayService) setSongPlaying(orderId uint64) {
	s.OrderSongRepo.UpdateStatus(orderId, 1) // Playing
}

func (s *MusicPlayService) resetOtherOrderSongStatus(liveRecordId uint64, playingOrderId uint64) {
	songs, _ := s.OrderSongRepo.SelectByLiveRecordId(liveRecordId)
	for _, song := range songs {
		if song.Status == 1 && song.Id != playingOrderId {
			s.OrderSongRepo.UpdateStatus(song.Id, 0) // Waiting
		}
	}
}

func (s *MusicPlayService) getOperatorInfoDto(userUuid string) *dto.BasicUserDto {
	user, _ := s.UserRepo.SelectByUserUuid(userUuid)
	if user != nil {
		return &dto.BasicUserDto{
			UserUuid: user.UserUuid,
			UserName: user.UserName,
			Icon:     user.Icon,
		}
	}
	return &dto.BasicUserDto{UserUuid: userUuid}
}

func (s *MusicPlayService) sendMusicPlayerMessage(event dto.MusicPlayNotifyEventDto, roomUuid string, eventType int) {
	type EventDto struct {
		Data      interface{} `json:"data"`
		EventType int         `json:"eventType"`
	}
	evt := EventDto{
		Data:      event,
		EventType: eventType,
	}
	msgBytes, _ := json.Marshal(evt)

	param := client.NeRoomMessageParam{
		RoomUuid: roomUuid,
		Message:  string(msgBytes),
	}
	s.NeRoomClient.SendNeRoomCustomMessage(param)
}

func (s *MusicPlayService) CleanPlayerMusicInfo(liveRecordId uint64) error {
	return s.MusicPlayerRepo.DeleteMusicPlayerInfo(liveRecordId)
}
