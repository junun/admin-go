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
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	SecretDing 		SecretDingTalkClient
	KeyWordDing 	KeyWordDingTalkClient
	AclDing			AclDingTalkClient
	WebChatApp 		WebChatAppClient
	WebChatRobot 	WebChatRobotClient
)

var DingList = []string{"1581xxxx723","1581xxxx724", "1581xxxx725"}
const DingUser = "1581xxxx723"

//func init() {
//	// 获取 dingsecret 配置
//	ds, err := setting.Cfg.GetSection("dingsecret")
//	if err != nil {
//		log.Fatal(2, "Fail to get section 'dingsecret': %v", err)
//	}
//
//	SecretDing = ReturnSecretDingTalkClient(
//		ds.Key("webhook").String(),
//		ds.Key("secret").String())
//
//	// 获取 dingkeyword 配置
//	dk, err := setting.Cfg.GetSection("dingkeyword")
//	if err != nil {
//		log.Fatal(2, "Fail to get section 'dingkeyword': %v", err)
//	}
//
//	KeyWordDing = ReturnKeyWordDingTalkClient(
//		dk.Key("webhook").String())
//
//
//	// 获取 dingacl 配置
//	da, err := setting.Cfg.GetSection("dingacl")
//	if err != nil {
//		log.Fatal(2, "Fail to get section 'dingacl': %v", err)
//	}
//
//	AclDing = ReturnAclDingTalkClient(
//		da.Key("webhook").String())
//
//}

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

type WebChatAppClient struct {
	webhook string
}

type WebChatRobotClient struct {
	webhook string
}

type Text struct {
	Content string `json:"content"`
}

type Markdown struct {
	Title string 	`json:"title"`
	Text  string 	`json:"text"`
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

type WebChatRobotSendRequest struct {
	MsgType  			string   `json:"msgtype"`
	Text     			Text     `json:"text"`
	Markdown 			Markdown `json:"markdown"`
	MentionedMobileList	[]string `json:"mentioned_mobile_list"`
}

type WebChatSendRequest struct {
	MsgType  string   `json:"msgtype"`
	Text     Text     `json:"text"`
	Agentid  string	  `json:"agentid"`
	Safe  	 string	  `json:"safe"`
	Toparty  string	  `json:"toparty"`
	Touser   string	  `json:"touser"`
	Totag    string	  `json:"totag"`
}

type OapiRobotSendResponse struct {
	ErrMsg  string `json:"errmsg"`
	ErrCode int64  `json:"errcode"`
}

type WebChatAccessTokenResponse struct {
	Errcode  	int  	`json:"errcode"`
	Errmsg  	string  `json:"errmsg"`
	AccessToken string  `json:"access_token"`
	ExpiresIn	int64 	`json:"expires_in"`
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

func ReturnWebChatAppClient(webhook string) WebChatAppClient {
	return WebChatAppClient{
		webhook: webhook,
	}
}

func ReturnWebChatRobotClient(webhook string) WebChatRobotClient {
	return WebChatRobotClient{
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

func CreateWebChatRobotSendTextRequest(content string,  mml []string, isAll bool)  WebChatRobotSendRequest {
	if isAll {
		mml = append(mml, "@all")
	}
	return WebChatRobotSendRequest {
		MsgType: "text",
		Text:    Text{Content: content},
		MentionedMobileList: mml,
	}
}

func CreateWebChatSendTextRequest(content, agentid, toparty, touser, totag string)  WebChatSendRequest{
	return WebChatSendRequest{
		MsgType: "text",
		Text:    Text{Content: content},
		Safe: "0",
		Agentid: agentid,
		Toparty: toparty,
		Touser: touser,
		Totag: totag,
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

func (d WebChatAppClient) Execute(request WebChatSendRequest) (*OapiRobotSendResponse, error) {
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

func (d WebChatRobotClient) Execute(request WebChatRobotSendRequest) (*OapiRobotSendResponse, error) {
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

func DingtalkSentTest(dtype int, webhook, secret, keyword string) error {
	switch dtype {
	case 1:
		if webhook != "" && secret != "" {
			SecretDing = ReturnSecretDingTalkClient(webhook, secret)
			textReq := CreateOapiRobotSendTextRequest(
				"钉钉数字签名测试消息",
				[]string{},
				false)
			if _, e := SecretDing.Execute(textReq); e != nil {
				return e
			}
		} else {
			return  errors.New("钉钉数字签名必须提供webhook地址和secret秘钥")
		}
	case 2:
		if webhook != "" && keyword != "" {
			KeyWordDing = ReturnKeyWordDingTalkClient(webhook)
			textReq := CreateOapiRobotSendTextRequest(
				keyword + "-钉钉关键字测试消息",
				[]string{},
				false)
			if _, e := KeyWordDing.Execute(textReq); e != nil {
				return e
			}
		} else {
			return  errors.New("钉钉关键字必须提供webhook地址和keyword关键词")
		}

	case 3:
		if webhook != "" {
			AclDing = ReturnAclDingTalkClient(webhook)
			textReq := CreateOapiRobotSendTextRequest(
				"钉钉Acl测试消息",
				[]string{},
				false)
			if _, e := AclDing.Execute(textReq); e != nil {
				return e
			}
		} else {
			return  errors.New("钉钉Acl必须提供webhook地址并设置放行的公网ip地址或者网段")
		}

	case 4:
		if webhook != "" && keyword != "" && secret != "" {
			accessToken, e	:= ReturnWebChatAccessToken(webhook, secret)
			if e != nil {
				return e
			}

			reqUrl := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + accessToken
			WebChatApp = ReturnWebChatAppClient(reqUrl)
			textReq	:= CreateWebChatSendTextRequest(
				"企业微信应用测试消息",
				keyword,
				"1",
				"","")

			if _, e := WebChatApp.Execute(textReq); e != nil {
				return e
			}
		} else {
			return  errors.New("企业微信应用必须提供公司id，应用id和应用秘钥")
		}

	case 5:
		if webhook != "" {
			WebChatRobot = ReturnWebChatRobotClient(webhook)
			textReq := CreateWebChatRobotSendTextRequest (
				"企业微信机器人测试",
				[]string{},
				false)
			if _, e := WebChatRobot.Execute(textReq); e != nil {
				return e
			}
		} else {
			return  errors.New("企业微信机器人必须提供webhook地址")
		}
	}

	return nil
}

func DingtalkSentChannel(id int, content string, atMobiles []string, isAtAll bool) error {
	strid := strconv.Itoa(id)
	msg := CreateOapiRobotSendTextRequest(content, atMobiles, isAtAll)
	if id == 0  {
		var webhook, secret string

		// 获取 dingsecret 配置
		ds, e := setting.Cfg.GetSection("dingsecret")
		if e != nil {
			return e
		}
		webhook = ds.Key("webhook").String()
		secret	= ds.Key("secret").String()

		SecretDing = ReturnSecretDingTalkClient(webhook, secret)
		if _, e := SecretDing.Execute(msg); e != nil {
			return e
		}
		return nil
	}

	var robot SettingRobot
	if e := DB.Find(&robot, strid).Error; e != nil {
		return e
	}

	switch robot.Type {
	case 1:
		SecretDing = ReturnSecretDingTalkClient(robot.Webhook, robot.Secret)
		if _, e := SecretDing.Execute(msg); e != nil {
			return e
		}
	case 2:
		KeyWordDing = ReturnKeyWordDingTalkClient(robot.Webhook)
		if _, e := KeyWordDing.Execute(msg); e != nil {
			return e
		}
	case 3:
		AclDing = ReturnAclDingTalkClient(robot.Webhook)
		if _, e := AclDing.Execute(msg); e != nil {
			return e
		}

	case 4:
		accessToken, e	:= ReturnWebChatAccessToken(robot.Webhook, robot.Secret)
		if e != nil {
			return e
		}
		reqUrl := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + accessToken

		WebChatApp 	= ReturnWebChatAppClient(reqUrl)
		textReq	:= CreateWebChatSendTextRequest (
			content,
			robot.Keyword,
			"1",
			"",
			"")

		if _, e := WebChatApp.Execute(textReq); e != nil {
			return e
		}

	case 5:
		WebChatRobot = ReturnWebChatRobotClient(robot.Webhook)
		textReq := CreateWebChatRobotSendTextRequest (
			content,
			atMobiles,
			isAtAll)
		if _, e := WebChatRobot.Execute(textReq); e != nil {
			return e
		}
	}

	return nil
}

func ReturnWebChatAccessToken(corpid, corpsecret string) (string, error) {
	key := WebChatAccessToken + corpid
	if Rdb.Exists(key).Val() == 1 {
		return  Rdb.Get(key).Val(), nil
	}

	requrl	:= "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpid + "&corpsecret=" + corpsecret
	resp, e := http.Get(requrl)
	if e 	!= nil {
		return "", e
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("http response status code is: %d", resp.StatusCode))
	}

	responseBody, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return "", e
	}

	var oResponse WebChatAccessTokenResponse
	if e = json.Unmarshal(responseBody, &oResponse); e != nil {
		return "", e
	}

	if oResponse.Errcode != 0 {
		return "", errors.New(fmt.Sprintf("response: %s", responseBody))
	}

	expired := time.Duration(oResponse.ExpiresIn - 200)
	Rdb.Set(key, oResponse.AccessToken, expired * time.Second)

	return oResponse.AccessToken, nil
}