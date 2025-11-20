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
	"net/url"
	"netease-kit/nemo/internal/config"
	"strconv"
	"time"
)

type NimClient struct {
	AppKey    string
	AppSecret string
	Host      string
}

func NewNimClient() *NimClient {
	return &NimClient{
		AppKey:    config.AppConfig.Yunxin.Origin.AppKey,
		AppSecret: config.AppConfig.Yunxin.Origin.AppSecret,
		Host:      config.AppConfig.Yunxin.Origin.NimHost,
	}
}

func (c *NimClient) CreateUser(accid, name, icon, token string) error {
	endpoint := "/nimserver/user/create.action"
	params := url.Values{}
	params.Set("accid", accid)
	params.Set("name", name)
	params.Set("icon", icon)
	params.Set("token", token)

	resp, err := c.post(endpoint, params)
	if err != nil {
		return err
	}

	if resp.Code != 200 {
		if resp.Desc == "already register" {
			return nil // Treat as success
		}
		return fmt.Errorf("nim create user failed: %s", resp.Desc)
	}
	return nil
}

func (c *NimClient) AddFriend(accid, faccid string, typeInt int, msg string) error {
	endpoint := "/nimserver/friend/add.action"
	params := url.Values{}
	params.Set("accid", accid)
	params.Set("faccid", faccid)
	params.Set("type", strconv.Itoa(typeInt))
	params.Set("msg", msg)

	resp, err := c.post(endpoint, params)
	if err != nil {
		return err
	}

	if resp.Code != 200 {
		return fmt.Errorf("nim add friend failed: %s", resp.Desc)
	}
	return nil
}

func (c *NimClient) SendImCustomMsg(from, to string, ope int, body interface{}) error {
	endpoint := "/nimserver/msg/sendMsg.action"
	params := url.Values{}
	params.Set("from", from)
	params.Set("ope", strconv.Itoa(ope))
	params.Set("to", to)
	params.Set("type", "100") // 100: Custom Message

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	params.Set("body", string(jsonBody))

	resp, err := c.post(endpoint, params)
	if err != nil {
		return err
	}

	if resp.Code != 200 {
		return fmt.Errorf("nim send custom msg failed: %s", resp.Desc)
	}
	return nil
}

func (c *NimClient) SendBatchAttachMsg(fromAccid string, toAccids []string, attach interface{}) error {
	endpoint := "/nimserver/msg/sendBatchAttachMsg.action"
	params := url.Values{}
	params.Set("fromAccid", fromAccid)

	toAccidsJson, err := json.Marshal(toAccids)
	if err != nil {
		return err
	}
	params.Set("toAccids", string(toAccidsJson))

	attachJson, err := json.Marshal(attach)
	if err != nil {
		return err
	}
	params.Set("attach", string(attachJson))

	resp, err := c.post(endpoint, params)
	if err != nil {
		return err
	}

	if resp.Code != 200 {
		return fmt.Errorf("nim send batch attach msg failed: %s", resp.Desc)
	}
	return nil
}

type NimResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func (c *NimClient) post(endpoint string, params url.Values) (*NimResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", c.Host+endpoint, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, err
	}

	c.addHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var nimResp NimResponse
	if err := json.Unmarshal(body, &nimResp); err != nil {
		return nil, err
	}
	return &nimResp, nil
}

func (c *NimClient) addHeaders(req *http.Request) {
	nonce := strconv.Itoa(rand.Intn(100000))
	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	checkSum := c.getCheckSum(nonce, curTime)

	req.Header.Set("AppKey", c.AppKey)
	req.Header.Set("Nonce", nonce)
	req.Header.Set("CurTime", curTime)
	req.Header.Set("CheckSum", checkSum)
}

func (c *NimClient) getCheckSum(nonce, curTime string) string {
	str := c.AppSecret + nonce + curTime
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
