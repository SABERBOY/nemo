package service

import (
	"errors"
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/model"
	"netease-kit/nemo/internal/repo"
	"time"
)

type GameService struct {
	GameRecordRepo *repo.GameRecordRepository
	GameMemberRepo *repo.GameMemberRepository
	UserRepo       *repo.UserRepository
}

func NewGameService() *GameService {
	return &GameService{
		GameRecordRepo: repo.NewGameRecordRepository(),
		GameMemberRepo: repo.NewGameMemberRepository(),
		UserRepo:       repo.NewUserRepository(),
	}
}

func (s *GameService) GetGameList() []dto.GameInfoDto {
	// Mock data
	return []dto.GameInfoDto{
		{GameId: 1, GameName: "Guess Song", GameDesc: "Guess the song name", Thumbnail: "url", Rule: "rule"},
	}
}

func (s *GameService) CreateGame(userUuid string, param dto.GameRoomParam) (*dto.GameRoomInfoDto, error) {
	// Check if game already exists
	exist, err := s.GameRecordRepo.SelectByRoomUuid(param.RoomUuid)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return s.toGameRoomInfoDto(exist), nil
	}

	record := &model.GameRecord{
		LiveRecordId:  param.LiveRecordId,
		RoomUuid:      param.RoomUuid,
		GameCreator:   userUuid,
		GameStatus:    0, // Waiting
		GameId:        param.GameId,
		GameName:      "Guess Song", // Mock
		RoomArchiveId: param.RoomArchiveId,
	}

	err = s.GameRecordRepo.Insert(record)
	if err != nil {
		return nil, err
	}

	return s.toGameRoomInfoDto(record), nil
}

func (s *GameService) JoinGame(userUuid string, param dto.GameRoomParam) (*dto.GameRoomMemberDto, error) {
	record, err := s.GameRecordRepo.SelectByRoomUuid(param.RoomUuid)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("game not found")
	}

	user, err := s.UserRepo.SelectByUserUuid(userUuid)
	if err != nil {
		return nil, err
	}

	member, err := s.GameMemberRepo.SelectByUserUuidAndGameRecordId(userUuid, record.Id)
	if err != nil {
		return nil, err
	}
	if member == nil {
		member = &model.GameMember{
			LiveRecordId:  record.LiveRecordId,
			RoomUuid:      record.RoomUuid,
			Status:        0, // Preparing
			JoinTime:      time.Now().UnixMilli(),
			GameId:        record.GameId,
			UserUuid:      userUuid,
			UserName:      user.UserName,
			RoomArchiveId: record.RoomArchiveId,
			GameRecordId:  record.Id,
		}
		err = s.GameMemberRepo.Insert(member)
		if err != nil {
			return nil, err
		}
	}

	return &dto.GameRoomMemberDto{
		UserUuid: userUuid,
		UserName: user.UserName,
		Avatar:   user.Icon,
		Status:   member.Status,
	}, nil
}

func (s *GameService) StartGame(userUuid string, param dto.GameRoomParam) error {
	record, err := s.GameRecordRepo.SelectByRoomUuid(param.RoomUuid)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.New("game not found")
	}
	if record.GameCreator != userUuid {
		return errors.New("permission denied")
	}

	return s.GameRecordRepo.UpdateStatus(record.Id, 1) // Playing
}

func (s *GameService) EndGame(userUuid string, param dto.GameRoomParam) error {
	record, err := s.GameRecordRepo.SelectByRoomUuid(param.RoomUuid)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.New("game not found")
	}
	if record.GameCreator != userUuid {
		return errors.New("permission denied")
	}

	return s.GameRecordRepo.UpdateStatus(record.Id, 2) // Ended
}

func (s *GameService) ExitGame(userUuid string, param dto.GameRoomParam) error {
	record, err := s.GameRecordRepo.SelectByRoomUuid(param.RoomUuid)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.New("game not found")
	}

	member, err := s.GameMemberRepo.SelectByUserUuidAndGameRecordId(userUuid, record.Id)
	if err != nil {
		return err
	}
	if member != nil {
		return s.GameMemberRepo.UpdateStatus(member.Id, 2) // Left
	}
	return nil
}

func (s *GameService) GetGameRoomMembers(userUuid string, param dto.GameRoomParam) ([]dto.GameRoomMemberDto, error) {
	record, err := s.GameRecordRepo.SelectByRoomUuid(param.RoomUuid)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return []dto.GameRoomMemberDto{}, nil
	}

	members, err := s.GameMemberRepo.SelectByGameRecordId(record.Id)
	if err != nil {
		return nil, err
	}

	var dtos []dto.GameRoomMemberDto
	for _, m := range members {
		dtos = append(dtos, dto.GameRoomMemberDto{
			UserUuid: m.UserUuid,
			UserName: m.UserName,
			Status:   m.Status,
		})
	}
	return dtos, nil
}

func (s *GameService) GetGameInfo(param dto.GameInfoParam) (*dto.GameRoomInfoDto, error) {
	record, err := s.GameRecordRepo.SelectByRoomUuid(param.RoomUuid)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("game not found")
	}
	return s.toGameRoomInfoDto(record), nil
}

func (s *GameService) StatusReporter(userUuid string, param dto.GameRoomParam) error {
	// Just a placeholder for now, as per Java implementation it might update status or heartbeat
	return nil
}

func (s *GameService) toGameRoomInfoDto(record *model.GameRecord) *dto.GameRoomInfoDto {
	return &dto.GameRoomInfoDto{
		GameId:        record.GameId,
		GameName:      record.GameName,
		GameDesc:      record.GameDesc,
		Thumbnail:     record.Thumbnail,
		Rule:          record.Rule,
		RoomUuid:      record.RoomUuid,
		GameCreator:   record.GameCreator,
		GameStatus:    record.GameStatus,
		GameRecordId:  record.Id,
		RoomArchiveId: record.RoomArchiveId,
	}
}
