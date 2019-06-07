package jssdk

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/uzicloud/gowechat/util"
	"github.com/uzicloud/gowechat/wxcontext"
)

const getTicketURL = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"

// Js struct
type Js struct {
	*wxcontext.Context
}

// Config 返回给用户jssdk配置信息
type Config struct {
	AppID     string `json:"app_id"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonce_str"`
	Signature string `json:"signature"`
}

//JsString wx.config中的配置
func (c *Config) ToMap() (cfg map[string]interface{}) {
	cfg = make(map[string]interface{})
	cfg["appId"] = c.AppID
	cfg["timestamp"] = c.Timestamp
	cfg["nonceStr"] = c.NonceStr
	cfg["signature"] = c.Signature
	return
}

// resTicket 请求jsapi_tikcet返回结果
type resTicket struct {
	util.CommonError

	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

//NewJs init
func NewJs(context *wxcontext.Context) *Js {
	js := new(Js)
	js.Context = context
	return js
}

//GetConfig 获取jssdk需要的配置参数
//uri 为当前网页地址
func (js *Js) GetConfig(url string) (config *Config, err error) {
	config = new(Config)
	var ticketStr string
	ticketStr, err = js.GetTicket()
	if err != nil {
		return
	}

	nonceStr := util.RandomStr(16)
	timestamp := util.GetCurrTs()
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticketStr, nonceStr, timestamp, url)
	sigStr := util.Signature(str)

	config.AppID = js.AppID
	config.NonceStr = nonceStr
	config.Timestamp = timestamp
	config.Signature = sigStr
	return
}

//GetTicket 获取jsapi_tocket
func (js *Js) GetTicket() (ticketStr string, err error) {
	js.GetJsAPITicketLock().Lock()
	defer js.GetJsAPITicketLock().Unlock()

	//先从cache中取
	jsAPITicketCacheKey := fmt.Sprintf("jsapi_ticket_%s", js.AppID)
	val := js.Cache.Get(jsAPITicketCacheKey)
	if val != nil {
		ticketStr = val.(string)
		return
	}
	var ticket resTicket
	ticket, err = js.getTicketFromServer()
	if err != nil {
		return
	}
	ticketStr = ticket.Ticket
	return
}

//getTicketFromServer 强制从服务器中获取ticket
func (js *Js) getTicketFromServer() (ticket resTicket, err error) {
	var accessToken string
	accessToken, err = js.GetAccessToken()
	if err != nil {
		return
	}

	var response []byte
	url := fmt.Sprintf(getTicketURL, accessToken)
	response, err = util.HTTPGet(url)
	err = json.Unmarshal(response, &ticket)
	if err != nil {
		return
	}
	if ticket.ErrCode != 0 {
		err = fmt.Errorf("getTicket Error : errcode=%d , errmsg=%s", ticket.ErrCode, ticket.ErrMsg)
		return
	}

	jsAPITicketCacheKey := fmt.Sprintf("jsapi_ticket_%s", js.AppID)
	expires := ticket.ExpiresIn - 1500
	err = js.Cache.Put(jsAPITicketCacheKey, ticket.Ticket, time.Duration(expires)*time.Second)
	return
}
