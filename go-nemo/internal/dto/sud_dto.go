package dto

type SudLoginDto struct {
	Code   string `json:"code"`
	AppId  string `json:"appId"`
	AppKey string `json:"appKey"`
}

type SudUserDto struct {
	Uid       string `json:"uid"`
	NickName  string `json:"nick_name"`
	AvatarUrl string `json:"avatar_url"`
	Gender    string `json:"gender"`
}

type GetSSTokenResp struct {
	SsToken    string      `json:"ss_token"`
	ExpireDate int64       `json:"expire_date"`
	UserInfo   *SudUserDto `json:"user_info"`
}

type UpdateSSTokenResp struct {
	SsToken    string `json:"ss_token"`
	ExpireDate int64  `json:"expire_date"`
}

type GetSSTokenParam struct {
	Code string `json:"code"`
}

type UpdateSSTokenParam struct {
	SsToken string `json:"ss_token"`
}

type GetUserInfoParam struct {
	SsToken string `json:"ss_token"`
}

type GameStartDto struct {
	ReportGameInfoKey string `json:"report_game_info_key"`
}

type GameSettleDto struct {
	ReportGameInfoKey string `json:"report_game_info_key"`
}
