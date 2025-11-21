package service

import (
	"netease-kit/nemo/internal/config"
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/repo"
	"netease-kit/nemo/pkg/utils"
)

type SudService struct {
	UserRepo *repo.UserRepository
}

func NewSudService() *SudService {
	return &SudService{
		UserRepo: repo.NewUserRepository(),
	}
}

func (s *SudService) SudLogin(appKey, userUuid string) (*dto.SudLoginDto, error) {
	code, err := s.GetCode(appKey, userUuid)
	if err != nil {
		return nil, err
	}
	return &dto.SudLoginDto{
		Code:   code,
		AppId:  config.AppConfig.Business.Game.AppId,
		AppKey: config.AppConfig.Business.Game.AppKey,
	}, nil
}

func (s *SudService) GetCode(appKey, userUuid string) (string, error) {
	return utils.CreateToken(appKey, userUuid, config.AppConfig.Business.Game.AppId, config.AppConfig.Business.Game.AppSecret)
}

func (s *SudService) GetSSToken(code string) (*dto.GetSSTokenResp, error) {
	claims, err := utils.VerifyToken(config.AppConfig.Business.Game.AppSecret, code)
	if err != nil {
		return nil, err
	}

	ssToken, expireDate, err := utils.CreateSSToken(claims.AppKey, claims.UserUuid, config.AppConfig.Business.Game.AppId, config.AppConfig.Business.Game.AppSecret, utils.DefaultMinSSTokenExpire)
	if err != nil {
		return nil, err
	}

	user, err := s.UserRepo.SelectByUserUuid(claims.UserUuid)
	if err != nil {
		return nil, err
	}

	sudUser := &dto.SudUserDto{
		Uid:       user.UserUuid,
		NickName:  user.UserName,
		AvatarUrl: user.Icon,
		Gender:    "male", // Default or map from user.Sex
	}
	if user.Sex == 2 {
		sudUser.Gender = "female"
	}

	return &dto.GetSSTokenResp{
		SsToken:    ssToken,
		ExpireDate: expireDate,
		UserInfo:   sudUser,
	}, nil
}

func (s *SudService) UpdateSSToken(ssToken string) (*dto.UpdateSSTokenResp, error) {
	claims, err := utils.VerifyToken(config.AppConfig.Business.Game.AppSecret, ssToken)
	if err != nil {
		return nil, err
	}

	newSSToken, expireDate, err := utils.CreateSSToken(claims.AppKey, claims.UserUuid, config.AppConfig.Business.Game.AppId, config.AppConfig.Business.Game.AppSecret, utils.DefaultMinSSTokenExpire)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateSSTokenResp{
		SsToken:    newSSToken,
		ExpireDate: expireDate,
	}, nil
}

func (s *SudService) GetUserInfo(ssToken string) (*dto.SudUserDto, error) {
	claims, err := utils.VerifyToken(config.AppConfig.Business.Game.AppSecret, ssToken)
	if err != nil {
		return nil, err
	}

	user, err := s.UserRepo.SelectByUserUuid(claims.UserUuid)
	if err != nil {
		return nil, err
	}

	sudUser := &dto.SudUserDto{
		Uid:       user.UserUuid,
		NickName:  user.UserName,
		AvatarUrl: user.Icon,
		Gender:    "male",
	}
	if user.Sex == 2 {
		sudUser.Gender = "female"
	}

	return sudUser, nil
}
