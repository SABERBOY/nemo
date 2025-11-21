package dto

type GameRoomParam struct {
	RoomUuid      string `json:"roomUuid"`
	RoomName      string `json:"roomName"`
	GameId        uint64 `json:"gameId"`
	LiveRecordId  uint64 `json:"liveRecordId"`
	RoomArchiveId string `json:"roomArchiveId"`
}

type GameInfoDto struct {
	GameId    uint64 `json:"gameId"`
	GameName  string `json:"gameName"`
	GameDesc  string `json:"gameDesc"`
	Thumbnail string `json:"thumbnail"`
	Rule      string `json:"rule"`
}

type GameRoomInfoDto struct {
	GameId        uint64 `json:"gameId"`
	GameName      string `json:"gameName"`
	GameDesc      string `json:"gameDesc"`
	Thumbnail     string `json:"thumbnail"`
	Rule          string `json:"rule"`
	RoomUuid      string `json:"roomUuid"`
	RoomName      string `json:"roomName"`
	GameCreator   string `json:"gameCreator"`
	GameStatus    int    `json:"gameStatus"`
	GameRecordId  uint64 `json:"gameRecordId"`
	RoomArchiveId string `json:"roomArchiveId"`
}

type GameRoomMemberDto struct {
	UserUuid string `json:"userUuid"`
	UserName string `json:"userName"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status"`
}

type GameInfoParam struct {
	RoomUuid string `form:"roomUuid"`
}
