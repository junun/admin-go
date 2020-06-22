package models

import (
	"api/pkg/setting"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	SecretDing 	SecretDingTalkClient
	KeyWordDing KeyWordDingTalkClient
	AclDing		AclDingTalkClient
)

var DingList = []string{"15818699723","15818699724", "15818699725"}
const DingUser = "15818699723"

func init() {
	// 获取 dingsecret 配置
	ds, err := setting.Cfg.GetSection("dingsecret")
	if err != nil {
		log.Fatal(2, "Fail to get section 'dingsecret': %v", err)
	}

	SecretDing = ReturnSecretDingTalkClient(
		ds.Key("webhook").String(),
		ds.Key("secret").String())

	// 获取 dingkeyword 配置
	dk, err := setting.Cfg.GetSection("dingkeyword")
	if err != nil {
		log.Fatal(2, "Fail to get section 'dingkeyword': %v", err)
	}

	KeyWordDing = ReturnKeyWordDingTalkClient(
		dk.Key("webhook").String())


	// 获取 dingacl 配置
	da, err := setting.Cfg.GetSection("dingacl")
	if err != nil {
		log.Fatal(2, "Fail to get section 'dingacl': %v", err)
	}

	AclDing = ReturnAclDingTalkClient(
		da.Key("webhook").String())

}


// Secret
type SecretDingTalkClient struct {
	webhook string
	secret  string
}

// 关键字
type KeyWordDingTalkClient struct {
	webhook string
}

// 地址白名单 支持两种设置方式：IP、IP段，暂不支持IPv6地址白名单
type AclDingTalkClient struct {
	webhook string
}

type Text struct {
	Content string `json:"content"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Link struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	PicUrl     string `json:"picUrl"`
	MessageUrl string `json:"messageUrl"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type OapiRobotSendRequest struct {
	MsgType  string   `json:"msgtype"`
	Text     Text     `json:"text"`
	Markdown Markdown `json:"markdown"`
	Link     Link     `json:"link"`
	At       At       `json:"at"`
}


type OapiRobotSendResponse struct {
	ErrMsg  string `json:"errmsg"`
	ErrCode int64  `json:"errcode"`
}

func ReturnSecretDingTalkClient(webhook, secret string) SecretDingTalkClient {
	return SecretDingTalkClient{
		webhook: webhook,
		secret:  secret,
	}
}

func ReturnKeyWordDingTalkClient(webhook string) KeyWordDingTalkClient {
	return KeyWordDingTalkClient{
		webhook: webhook,
	}
}

func ReturnAclDingTalkClient(webhook string) AclDingTalkClient {
	return AclDingTalkClient{
		webhook: webhook,
	}
}

func CreateOapiRobotSendTextRequest(content string, atMobiles []string, isAtAll bool) OapiRobotSendRequest {
	return OapiRobotSendRequest{
		MsgType: "text",
		Text:    Text{Content: content},
		At:      At{AtMobiles: atMobiles, IsAtAll: isAtAll},
	}

}

func (d SecretDingTalkClient) Execute(request OapiRobotSendRequest) (*OapiRobotSendResponse, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	pushUrl, err := d.getPushUrl()
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(pushUrl, "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("http response status code is: %d", resp.StatusCode))
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var oResponse OapiRobotSendResponse
	if err = json.Unmarshal(responseBody, &oResponse); err != nil {
		return nil, err
	}
	if oResponse.ErrCode != 0 {
		return &oResponse, errors.New(fmt.Sprintf("response: %s", responseBody))
	}
	return &oResponse, nil
}

func (d SecretDingTalkClient) getPushUrl() (string, error) {
	if d.secret == "" {
		return d.webhook, nil
	}

	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	sign, err := d.sign(timestamp)
	if err != nil {
		return d.webhook, err
	}

	query := url.Values{}
	query.Set("timestamp", timestamp)
	query.Set("sign", sign)
	return d.webhook + "&" + query.Encode(), nil
}

func (d SecretDingTalkClient) sign(timestamp string) (string, error) {
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, d.secret)
	h := hmac.New(sha256.New, []byte(d.secret))
	if _, err := io.WriteString(h, stringToSign); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func (d KeyWordDingTalkClient) Execute(request OapiRobotSendRequest) (*OapiRobotSendResponse, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(d.webhook, "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("http response status code is: %d", resp.StatusCode))
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var oResponse OapiRobotSendResponse
	if err = json.Unmarshal(responseBody, &oResponse); err != nil {
		return nil, err
	}
	if oResponse.ErrCode != 0 {
		return &oResponse, errors.New(fmt.Sprintf("response: %s", responseBody))
	}
	return &oResponse, nil
}

func (d AclDingTalkClient) Execute(request OapiRobotSendRequest) (*OapiRobotSendResponse, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(d.webhook, "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("http response status code is: %d", resp.StatusCode))
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var oResponse OapiRobotSendResponse
	if err = json.Unmarshal(responseBody, &oResponse); err != nil {
		return nil, err
	}
	if oResponse.ErrCode != 0 {
		return &oResponse, errors.New(fmt.Sprintf("response: %s", responseBody))
	}
	return &oResponse, nil
}
