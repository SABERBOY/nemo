package service

import (
	"errors"
	"netease-kit/nemo/internal/client"
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/model"
	"netease-kit/nemo/internal/repo"

	"github.com/google/uuid"
)

type KtvService struct {
	OrderSongRepo    *repo.OrderSongRepository
	ChorusRecordRepo *repo.ChorusRecordRepository
	LiveRecordRepo   *repo.LiveRecordRepository
	UserRepo         *repo.UserRepository
	NeRoomClient     *client.NeRoomClient
}

func NewKtvService() *KtvService {
	return &KtvService{
		OrderSongRepo:    repo.NewOrderSongRepository(),
		ChorusRecordRepo: repo.NewChorusRecordRepository(),
		LiveRecordRepo:   repo.NewLiveRecordRepository(),
		UserRepo:         repo.NewUserRepository(),
		NeRoomClient:     client.NewNeRoomClient(),
	}
}

// --- Order Song Service ---

func (s *KtvService) OrderSong(userUuid string, param dto.OrderSongParam) (*dto.OrderSongDto, error) {
	liveRecord, err := s.LiveRecordRepo.SelectById(param.LiveRecordId)
	if err != nil || liveRecord == nil {
		return nil, errors.New("live record not found")
	}

	orderSong := &model.OrderSong{
		LiveRecordId:  param.LiveRecordId,
		RoomArchiveId: liveRecord.RoomArchiveId,
		RoomUuid:      liveRecord.RoomUuid,
		UserUuid:      userUuid,
		SongId:        param.SongId,
		SongName:      param.SongName,
		SongCover:     param.SongCover,
		Singer:        param.Singer,
		SingerCover:   param.SingerCover,
		SongTime:      param.SongTime,
		Channel:       param.Channel,
		Status:        0, // Waiting
	}

	err = s.OrderSongRepo.Insert(orderSong)
	if err != nil {
		return nil, err
	}

	// TODO: Send NeRoom message about new song order

	return s.buildOrderSongDto(orderSong), nil
}

func (s *KtvService) GetOrderSongs(liveRecordId uint64) ([]*dto.OrderSongDto, error) {
	songs, err := s.OrderSongRepo.SelectByLiveRecordId(liveRecordId)
	if err != nil {
		return nil, err
	}

	var dtos []*dto.OrderSongDto
	for _, song := range songs {
		dtos = append(dtos, s.buildOrderSongDto(song))
	}
	return dtos, nil
}

func (s *KtvService) CancelOrderSong(userUuid string, liveRecordId uint64, orderId uint64) error {
	song, err := s.OrderSongRepo.SelectById(orderId)
	if err != nil || song == nil {
		return errors.New("song not found")
	}
	if song.UserUuid != userUuid {
		// Check if user is host
		liveRecord, _ := s.LiveRecordRepo.SelectById(liveRecordId)
		if liveRecord == nil || liveRecord.UserUuid != userUuid {
			return errors.New("permission denied")
		}
	}

	return s.OrderSongRepo.Delete(orderId)
}

func (s *KtvService) SongSetTop(userUuid string, liveRecordId uint64, orderId uint64) error {
	// Only host or song owner? Usually host or paid feature. Assuming host for now.
	liveRecord, err := s.LiveRecordRepo.SelectById(liveRecordId)
	if err != nil || liveRecord == nil {
		return errors.New("live record not found")
	}
	if liveRecord.UserUuid != userUuid {
		return errors.New("permission denied")
	}

	return s.OrderSongRepo.SetTop(orderId)
}

func (s *KtvService) CleanUserOrderSongs(liveRecordId uint64, userUuid string) error {
	// Logic to clean all songs by a user.
	// For simplicity, just fetch and delete loop or add repo method.
	// Skipping for now or implementing simple loop.
	songs, _ := s.OrderSongRepo.SelectByLiveRecordId(liveRecordId)
	for _, song := range songs {
		if song.UserUuid == userUuid {
			s.OrderSongRepo.Delete(song.Id)
		}
	}
	return nil
}

func (s *KtvService) SwitchSong(liveRecordId uint64, userUuid string, currentOrderId uint64, attachment string) error {
	// Mark current as played
	if currentOrderId > 0 {
		s.OrderSongRepo.UpdateStatus(currentOrderId, 2) // Played
	}

	// Logic to start next song?
	// Usually client handles "Start Sing" for next song.
	return nil
}

// --- Sing Service ---

func (s *KtvService) SingStart(userUuid string, param dto.SingParam) error {
	// Update song status to Playing
	err := s.OrderSongRepo.UpdateStatus(param.OrderId, 1)
	if err != nil {
		return err
	}
	return nil
}

func (s *KtvService) GetSingInfo(roomUuid string) (*dto.SingDetailInfoDto, error) {
	// Find live record by roomUuid
	// This is tricky without liveRecordId. Assuming we can find it or pass it.
	// For now, return empty or mock.
	return &dto.SingDetailInfoDto{}, nil
}

func (s *KtvService) SingControl(userUuid string, param dto.SingActionParam) error {
	// Handle Pause/Resume/End logic
	// Mostly messaging to room.
	return nil
}

func (s *KtvService) AbandonSing(userUuid string, param dto.AbandonSingParam) error {
	return s.OrderSongRepo.UpdateStatus(param.OrderId, -1) // Cancelled/Abandoned
}

// --- Chorus Service ---

func (s *KtvService) ChorusInvite(userUuid string, param dto.ChorusInviteParam) (*dto.ChorusControlResultDto, error) {
	// Create Chorus Record
	chorusId := uuid.New().String()
	record := &model.ChorusRecord{
		ChorusId:      chorusId,
		RoomUuid:      param.RoomUuid,
		OrderId:       param.OrderId,
		LeaderUuid:    userUuid,
		AssistantUuid: param.UserUuid, // Invitee
		Status:        1,
		State:         0, // Inviting
	}
	err := s.ChorusRecordRepo.Insert(record)
	if err != nil {
		return nil, err
	}

	return &dto.ChorusControlResultDto{Code: 200, Msg: "Invited"}, nil
}

func (s *KtvService) JoinChorus(userUuid string, param dto.JoinChorusParam) (*dto.ChorusControlResultDto, error) {
	record, err := s.ChorusRecordRepo.SelectByChorusId(param.ChorusId)
	if err != nil || record == nil {
		return nil, errors.New("chorus record not found")
	}

	// Update state to Agreed/Singing
	s.ChorusRecordRepo.UpdateState(record.Id, 5) // Agreed

	return &dto.ChorusControlResultDto{Code: 200, Msg: "Joined"}, nil
}

func (s *KtvService) CancelChorus(userUuid string, param dto.CancelChorusParam) (*dto.ChorusControlResultDto, error) {
	record, err := s.ChorusRecordRepo.SelectByChorusId(param.ChorusId)
	if err != nil || record == nil {
		return nil, errors.New("chorus record not found")
	}
	s.ChorusRecordRepo.UpdateState(record.Id, 3) // Cancelled
	return &dto.ChorusControlResultDto{Code: 200, Msg: "Cancelled"}, nil
}

func (s *KtvService) ChorusReady(userUuid string, param dto.FinishChorusReadyParam) (*dto.ChorusControlResultDto, error) {
	// Mark as ready
	return &dto.ChorusControlResultDto{Code: 200, Msg: "Ready"}, nil
}

func (s *KtvService) CleanOrderSongs(liveRecordId uint64) error {
	songs, _ := s.OrderSongRepo.SelectByLiveRecordId(liveRecordId)
	for _, song := range songs {
		s.OrderSongRepo.Delete(song.Id)
	}
	return nil
}

func (s *KtvService) CleanSingInfo(roomUuid string) error {
	// TODO: Implement Redis cleanup for sing info
	return nil
}

func (s *KtvService) EndSingWhenMemberOut(roomUuid string, userUuid string, oldOrderSongs []*dto.OrderSongDto) error {
	// TODO: Implement logic to handle singing user leaving
	return nil
}

// Helper
func (s *KtvService) buildOrderSongDto(song *model.OrderSong) *dto.OrderSongDto {
	user, _ := s.UserRepo.SelectByUserUuid(song.UserUuid)
	userName := ""
	userIcon := ""
	if user != nil {
		userName = user.UserName
		userIcon = user.Icon
	}

	return &dto.OrderSongDto{
		Id:            song.Id,
		LiveRecordId:  song.LiveRecordId,
		RoomArchiveId: song.RoomArchiveId,
		RoomUuid:      song.RoomUuid,
		UserUuid:      song.UserUuid,
		UserName:      userName,
		UserIcon:      userIcon,
		SongId:        song.SongId,
		SongName:      song.SongName,
		SongCover:     song.SongCover,
		Singer:        song.Singer,
		SingerCover:   song.SingerCover,
		SongTime:      song.SongTime,
		Channel:       song.Channel,
		Status:        song.Status,
		SetTopTime:    song.SetTopTime,
	}
}
