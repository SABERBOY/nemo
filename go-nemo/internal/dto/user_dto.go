package dto

type InitUserParam struct {
	UserName  string `json:"userName"`
	SceneType int    `json:"sceneType"` // 1: 1V1, 2: Voice Room
	UserUuid  string `json:"userUuid"`
	ImToken   string `json:"imToken"`
	UserToken string `json:"userToken"`
	Icon      string `json:"icon"`
}

type UserDto struct {
	UserUuid  string `json:"userUuid"`
	UserName  string `json:"userName"`
	Icon      string `json:"icon"`
	UserToken string `json:"userToken"`
	ImToken   string `json:"imToken"`
	Sex       int    `json:"sex"`
	Mobile    string `json:"mobile"`
}
