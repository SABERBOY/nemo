package service

import (
	"errors"
	"fmt"
	"math/rand"
	"netease-kit/nemo/internal/client"
	"netease-kit/nemo/internal/config"
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/model"
	"netease-kit/nemo/internal/repo"
	"time"

	"github.com/google/uuid"
)

type AppInitService struct {
	UserRepo     *repo.UserRepository
	NimClient    *client.NimClient
	NeRoomClient *client.NeRoomClient
}

func NewAppInitService() *AppInitService {
	return &AppInitService{
		UserRepo:     repo.NewUserRepository(),
		NimClient:    client.NewNimClient(),
		NeRoomClient: client.NewNeRoomClient(),
	}
}

func (s *AppInitService) InitAppAndUser(appKey string, param dto.InitUserParam) (*dto.UserDto, error) {
	if appKey == "" {
		return nil, errors.New("appKey cannot be empty")
	}

	// Initialize System Users (Logic simplified for migration: check if exists, if not create)
	s.initSystemUser()

	// Build User
	user := s.buildUser(param)

	// Init User based on SceneType
	return s.initUser(user, param.SceneType)
}

func (s *AppInitService) initSystemUser() {
	systemAccid := config.AppConfig.Business.SystemAccid
	assistAccid := config.AppConfig.Business.YunxinAssistAccid

	s.ensureUserExists(systemAccid, "System Bot")
	s.ensureUserExists(assistAccid, "Yunxin Assist")
}

func (s *AppInitService) ensureUserExists(userUuid string, name string) {
	user, _ := s.UserRepo.SelectByUserUuid(userUuid)
	if user == nil {
		token := uuid.New().String()
		newUser := &model.User{
			UserUuid:  userUuid,
			UserName:  name,
			UserToken: token,
			ImToken:   token,
			State:     1,
		}
		s.UserRepo.Insert(newUser)
		// Call NimService to create user on Yunxin
		s.NimClient.CreateUser(userUuid, name, "", token)
	}
}

func (s *AppInitService) buildUser(param dto.InitUserParam) *model.User {
	user := &model.User{
		UserName:  param.UserName,
		UserUuid:  param.UserUuid,
		Icon:      param.Icon,
		ImToken:   param.ImToken,
		UserToken: param.UserToken,
		Age:       -1,
		Sex:       -1,
		State:     1,
		Mobile:    fmt.Sprintf("%d", rand.Int63n(10000000000)), // Simplified random mobile
	}

	if user.UserName == "" {
		user.UserName = fmt.Sprintf("User_%d", time.Now().Unix())
	}
	if user.UserUuid == "" {
		user.UserUuid = uuid.New().String() // Simplified UUID generation
	}
	if user.Icon == "" {
		user.Icon = "default_icon_url" // Placeholder
	}
	token := uuid.New().String()
	if user.ImToken == "" {
		user.ImToken = token
	}
	if user.UserToken == "" {
		user.UserToken = token
	}

	return user
}

func (s *AppInitService) initUser(user *model.User, sceneType int) (*dto.UserDto, error) {
	existUser, err := s.UserRepo.SelectByUserUuid(user.UserUuid)
	if err != nil {
		return nil, err
	}
	if existUser != nil {
		return s.toUserDto(existUser), nil
	}

	// Insert User
	err = s.UserRepo.Insert(user)
	if err != nil {
		return nil, err
	}

	// Call NimService / NeRoomService based on sceneType
	if sceneType == 1 {
		// 1V1: Init IM User & Add Friend
		err = s.NimClient.CreateUser(user.UserUuid, user.UserName, user.Icon, user.ImToken)
		if err != nil {
			return nil, err
		}
		assistAccid := config.AppConfig.Business.YunxinAssistAccid
		s.NimClient.AddFriend(user.UserUuid, assistAccid, 1, "Add Yunxin Assist")
	} else if sceneType == 2 {
		// Voice Room: Init NeRoom User & Add Friend
		err = s.NeRoomClient.CreateNeRoomUser(user.UserUuid, client.CreateNeRoomUserParam{
			UserName:         user.UserName,
			Icon:             user.Icon,
			UserToken:        user.UserToken,
			ImToken:          user.ImToken,
			UpdateOnConflict: true,
		})
		if err != nil {
			return nil, err
		}
		assistAccid := config.AppConfig.Business.YunxinAssistAccid
		s.NimClient.AddFriend(user.UserUuid, assistAccid, 1, "Add Yunxin Assist")
	}

	return s.toUserDto(user), nil
}

func (s *AppInitService) toUserDto(user *model.User) *dto.UserDto {
	return &dto.UserDto{
		UserUuid:  user.UserUuid,
		UserName:  user.UserName,
		Icon:      user.Icon,
		UserToken: user.UserToken,
		ImToken:   user.ImToken,
		Sex:       user.Sex,
		Mobile:    user.Mobile,
	}
}
