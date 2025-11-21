package dto

type RoomMember struct {
	Role     string `json:"role"`
	UserUuid string `json:"userUuid"`
}

type CreateRoomEventNotify struct {
	RoomArchiveId string `json:"roomArchiveId"`
	RoomUuid      string `json:"roomUuid"`
}

type CloseRoomEventNotify struct {
	RoomArchiveId string `json:"roomArchiveId"`
	RoomUuid      string `json:"roomUuid"`
}

type JoinRoomEventNotify struct {
	RoomArchiveId string       `json:"roomArchiveId"`
	RoomUuid      string       `json:"roomUuid"`
	Users         []RoomMember `json:"users"`
}

type LeaveRoomEventNotify struct {
	RoomArchiveId string       `json:"roomArchiveId"`
	RoomUuid      string       `json:"roomUuid"`
	Users         []RoomMember `json:"users"`
}

type SeatUser struct {
	UserUuid        string `json:"user_uuid"`
	Name            string `json:"name"`
	Icon            string `json:"icon"`
	InviterUserUuid string `json:"inviter_user_uuid"`
	OnSeatType      int    `json:"on_seat_type"`
	Ext             string `json:"ext"`
	Timestamp       int64  `json:"timestamp"`
	Index           int    `json:"index"`
}

type UserOnSeatNotifyDto struct {
	RoomArchiveId string   `json:"room_archive_id"`
	UserUuid      string   `json:"user_uuid"`
	Index         int      `json:"index"`
	User          SeatUser `json:"user"`
}

type UserOffSeatNotifyDto struct {
	RoomArchiveId string `json:"room_archive_id"`
	UserUuid      string `json:"user_uuid"`
	Index         int    `json:"index"`
}

type NeRoomNotifyBody struct {
	EventType string      `json:"eventType"`
	Body      interface{} `json:"body"`
}
