package client

type CreateNeRoomParam struct {
	RoomName           string              `json:"roomName"`
	RoomUuid           string              `json:"roomUuid"`
	TemplateId         int64               `json:"templateId"`
	RoomSeatConfig     *RoomSeatConfig     `json:"roomSeatConfig,omitempty"`
	RoomConfig         *RoomConfig         `json:"roomConfig,omitempty"`
	RoomProfile        int                 `json:"roomProfile,omitempty"`
	Ext                string              `json:"ext,omitempty"`
	ExternalLiveConfig *ExternalLiveConfig `json:"externalLiveConfig,omitempty"`
}

type RoomSeatConfig struct {
	SeatCount  int `json:"seatCount"`
	ApplyMode  int `json:"applyMode"`
	InviteMode int `json:"inviteMode"`
}

type RoomConfig struct {
	Resource *ResourceConfig `json:"resource"`
}

type ResourceConfig struct {
	Rtc        bool `json:"rtc"`
	Chatroom   bool `json:"chatroom"`
	Whiteboard bool `json:"whiteboard"`
	Live       bool `json:"live"`
}

type ExternalLiveConfig struct {
	PushUrl     string `json:"pushUrl"`
	RtmpPullUrl string `json:"rtmpPullUrl"`
	HlsPullUrl  string `json:"hlsPullUrl"`
	HttpPullUrl string `json:"httpPullUrl"`
}

type CreateNeRoomDto struct {
	RoomArchiveId string `json:"roomArchiveId"`
	RoomUuid      string `json:"roomUuid"`
	Name          string `json:"name"`
}

type NeRoomSeatDto struct {
	User   *SeatUser `json:"user"`
	Status int       `json:"status"`
	Index  int       `json:"index"`
}

type SeatUser struct {
	UserUuid string `json:"userUuid"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
}

type NeRoomMemberDto struct {
	Total int64         `json:"total"`
	List  []*NeRoomUser `json:"list"`
}

type NeRoomUser struct {
	UserUuid string `json:"userUuid"`
	UserName string `json:"userName"`
	Role     string `json:"role"`
	Icon     string `json:"icon"`
}

type NeRoomMessageParam struct {
	RoomUuid string `json:"roomUuid"`
	Message  string `json:"message"` // JSON string
}

type CreateNeRoomParamV3 struct {
	RoomUuid     string                 `json:"roomUuid"`
	RoomName     string                 `json:"roomName"`
	TemplateId   int64                  `json:"templateId"`
	HostUserUuid string                 `json:"hostUserUuid"`
	RoomProfile  int                    `json:"roomProfile"`
	Config       *RoomComponentConfigV3 `json:"config"`
}

type RoomComponentConfigV3 struct {
	Seat     *SeatConfigV3 `json:"seat,omitempty"`
	Chatroom bool          `json:"chatroom"`
	Rtc      bool          `json:"rtc"`
	Live     *LiveConfigV3 `json:"live,omitempty"`
}

type SeatConfigV3 struct {
	Enable     bool `json:"enable"`
	SeatCount  int  `json:"seatCount"`
	ApplyMode  int  `json:"applyMode"`
	InviteMode int  `json:"inviteMode"`
}

type LiveConfigV3 struct {
	Enable bool `json:"enable"`
}

type StartLiveParam struct {
	RoomArchiveId string   `json:"roomArchiveId"`
	LiveUsers     []string `json:"liveUsers"`
	Topic         string   `json:"topic"`
}

type StartLiveResponseDto struct {
	TaskId string `json:"taskId"`
}
