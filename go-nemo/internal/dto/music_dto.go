package dto

type MusicActionParam struct {
	LiveRecordId uint64 `json:"liveRecordId"`
	Action       int    `json:"action"` // 1: Play, 2: Pause, 3: Resume, 4: End
	FirstPlay    bool   `json:"firstPlay"`
}

type MusicReadyParam struct {
	LiveRecordId uint64 `json:"liveRecordId"`
	OrderId      uint64 `json:"orderId"`
}

type PlayDetailInfoDto struct {
	OrderId     uint64 `json:"orderId"`
	SongId      string `json:"songId"`
	SongName    string `json:"songName"`
	SongCover   string `json:"songCover"`
	Singer      string `json:"singer"`
	SingerCover string `json:"singerCover"`
	SongTime    int64  `json:"songTime"`
	Channel     int    `json:"channel"`
	MusicStatus int    `json:"musicStatus"` // 0: Readying, 1: Play, 2: Pause
	RoomUuid    string `json:"roomUuid"`
	UserUuid    string `json:"userUuid"`
	UserName    string `json:"userName"`
	UserIcon    string `json:"userIcon"`
}

type MusicPlayNotifyEventDto struct {
	Operator      *BasicUserDto      `json:"operator"`
	PlayMusicInfo *PlayDetailInfoDto `json:"playMusicInfo"`
}

type BasicUserDto struct {
	UserUuid string `json:"userUuid"`
	UserName string `json:"userName"`
	Icon     string `json:"icon"`
}
