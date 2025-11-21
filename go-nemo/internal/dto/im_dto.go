package dto

type ImMemberInOutMsgCallbackParam struct {
	EventType  string `json:"eventType"`
	RoomId     string `json:"roomId"`
	Accid      string `json:"accid"`
	Event      string `json:"event"`
	RoleType   string `json:"roleType"`
	ClientType string `json:"clientType"`
	Code       string `json:"code"`
	ClientIp   string `json:"clientIp"`
	SdkVersion string `json:"sdkVersion"`
	Timestamp  int64  `json:"timestamp"`
	Tags       string `json:"tags"`
}
