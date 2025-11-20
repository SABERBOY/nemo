package client

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"netease-kit/nemo/internal/config"
	"strconv"
	"time"
)

type NeRoomClient struct {
	AppKey    string
	AppSecret string
	Host      string
}

func NewNeRoomClient() *NeRoomClient {
	return &NeRoomClient{
		AppKey:    config.AppConfig.Yunxin.Origin.AppKey,
		AppSecret: config.AppConfig.Yunxin.Origin.AppSecret,
		Host:      config.AppConfig.Yunxin.Origin.NeRoomHost,
	}
}

type CreateNeRoomUserParam struct {
	UserName         string `json:"userName"`
	Icon             string `json:"icon"`
	UserToken        string `json:"userToken"`
	ImToken          string `json:"imToken"`
	UpdateOnConflict bool   `json:"updateOnConflict"`
}

type NeRoomResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (c *NeRoomClient) CreateNeRoomUser(userUuid string, param CreateNeRoomUserParam) error {
	endpoint := fmt.Sprintf("/apps/%s/v1/users/%s", c.AppKey, userUuid)

	resp, err := c.request("PUT", endpoint, param)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return fmt.Errorf("neroom create user failed: %s", resp.Msg)
	}
	return nil
}

func (c *NeRoomClient) CreateNeRoom(param CreateNeRoomParam) (*CreateNeRoomDto, error) {
	endpoint := "/apps/v2/room"
	resp, err := c.request("PUT", endpoint, param)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("neroom create room failed: %s", resp.Msg)
	}

	var dto CreateNeRoomDto
	if err := json.Unmarshal(resp.Data, &dto); err != nil {
		return nil, err
	}
	return &dto, nil
}

func (c *NeRoomClient) DeleteNeRoom(roomArchiveId string) error {
	endpoint := fmt.Sprintf("/apps/v2/room?roomArchiveId=%s", roomArchiveId)
	resp, err := c.request("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf("neroom delete room failed: %s", resp.Msg)
	}
	return nil
}

func (c *NeRoomClient) GetNeRoomSeatList(roomArchiveId string) ([]*NeRoomSeatDto, error) {
	endpoint := fmt.Sprintf("/neroom/v1/seats/slot_list/%s", roomArchiveId)
	resp, err := c.request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("neroom get seat list failed: %s", resp.Msg)
	}

	var list []*NeRoomSeatDto
	if err := json.Unmarshal(resp.Data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (c *NeRoomClient) GetNeRoomOnlineMember(roomArchiveId string, pageNumber, pageSize int) (*NeRoomMemberDto, error) {
	endpoint := fmt.Sprintf("/apps/v2/online-user-list?roomArchiveId=%s&pageNumber=%d&pageSize=%d", roomArchiveId, pageNumber, pageSize)
	resp, err := c.request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("neroom get online member failed: %s", resp.Msg)
	}

	var dto NeRoomMemberDto
	if err := json.Unmarshal(resp.Data, &dto); err != nil {
		return nil, err
	}
	return &dto, nil
}

func (c *NeRoomClient) SendNeRoomCustomMessage(param NeRoomMessageParam) error {
	endpoint := "/apps/v2/sendChatRoomCustomMessage"
	resp, err := c.request("POST", endpoint, param)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf("neroom send custom message failed: %s", resp.Msg)
	}
	return nil
}

func (c *NeRoomClient) GetRoomMemberInfo(roomArchiveId, userUuid string) (*NeRoomUser, error) {
	endpoint := fmt.Sprintf("/apps/v2/room-user?roomArchiveId=%s&userUuid=%s", roomArchiveId, userUuid)
	resp, err := c.request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("neroom get room member info failed: %s", resp.Msg)
	}

	var user NeRoomUser
	if err := json.Unmarshal(resp.Data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *NeRoomClient) CreateNeRoomV3(param CreateNeRoomParamV3) (*CreateNeRoomDto, error) {
	endpoint := "/neroom/v4/rooms"
	resp, err := c.request("POST", endpoint, param)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("neroom create room v3 failed: %s", resp.Msg)
	}

	var dto CreateNeRoomDto
	if err := json.Unmarshal(resp.Data, &dto); err != nil {
		return nil, err
	}
	return &dto, nil
}

func (c *NeRoomClient) StartLive(param StartLiveParam) (*StartLiveResponseDto, error) {
	endpoint := "/apps/v2/live-start"
	resp, err := c.request("POST", endpoint, param)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("neroom start live failed: %s", resp.Msg)
	}

	var dto StartLiveResponseDto
	if err := json.Unmarshal(resp.Data, &dto); err != nil {
		return nil, err
	}
	return &dto, nil
}

func (c *NeRoomClient) UpdateLive(param StartLiveParam) (*StartLiveResponseDto, error) {
	endpoint := "/apps/v2/live-update"
	resp, err := c.request("POST", endpoint, param)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("neroom update live failed: %s", resp.Msg)
	}

	var dto StartLiveResponseDto
	if err := json.Unmarshal(resp.Data, &dto); err != nil {
		return nil, err
	}
	return &dto, nil
}

func (c *NeRoomClient) StopLive(roomArchiveId string) error {
	endpoint := "/apps/v2/live-end"
	body := map[string]string{"roomArchiveId": roomArchiveId}
	resp, err := c.request("POST", endpoint, body)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf("neroom stop live failed: %s", resp.Msg)
	}
	return nil
}

func (c *NeRoomClient) request(method, endpoint string, body interface{}) (*NeRoomResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.Host+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	c.addHeaders(req)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := ioutil.ReadAll(resp.Body)
	var neRoomResp NeRoomResponse
	if err := json.Unmarshal(respBytes, &neRoomResp); err != nil {
		return nil, err
	}
	return &neRoomResp, nil
}

func (c *NeRoomClient) addHeaders(req *http.Request) {
	nonce := strconv.Itoa(rand.Intn(100000))
	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	checkSum := c.getCheckSum(nonce, curTime)

	req.Header.Set("AppKey", c.AppKey)
	req.Header.Set("Nonce", nonce)
	req.Header.Set("CurTime", curTime)
	req.Header.Set("CheckSum", checkSum)
}

func (c *NeRoomClient) getCheckSum(nonce, curTime string) string {
	str := c.AppSecret + nonce + curTime
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
