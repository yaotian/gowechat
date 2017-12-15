package account

import (
	"encoding/json"
	"errors"

	"github.com/chanxuehong/wechat/mp"
	"github.com/yaotian/gowechat/mp/base"
	"github.com/yaotian/gowechat/wxcontext"
)

const (
	qrcodeURL = "https://api.weixin.qq.com/cgi-bin/qrcode/create"
)

//Qrcode 带参数的二维码
type Qrcode struct {
	base.MpBase
}

//NewQrcode 实例化
func NewQrcode(context *wxcontext.Context) *Qrcode {
	qrcode := new(Qrcode)
	qrcode.Context = context
	return qrcode
}

const (
	//TemporaryQRCodeExpireSecondsLimit 临时二维码 expire_seconds 限制
	TemporaryQRCodeExpireSecondsLimit = 604800
	//PermanentQRCodeSceneIDLimit 永久二维码 scene_id 限制
	PermanentQRCodeSceneIDLimit = 100000
)

//PermanentQRCode 永久二维码
type PermanentQRCode struct {
	// 下面两个字段同时只有一个有效, 非zero值表示有效.
	SceneID     uint32 `json:"scene_id,omitempty"`  // 场景值ID, 临时二维码时为32位非0整型, 永久二维码时最大值为100000(目前参数只支持1--100000)
	SceneString string `json:"scene_str,omitempty"` // 场景值ID(字符串形式的ID), 字符串类型, 长度限制为1到64, 仅永久二维码支持此字段

	Ticket string `json:"ticket"` // 获取的二维码ticket, 凭借此ticket可以在有效时间内换取二维码.
	URL    string `json:"url"`    // 二维码图片解析后的地址, 开发者可根据该地址自行生成需要的二维码图片
}

//TemporaryQRCode 临时二维码
type TemporaryQRCode struct {
	ExpireSeconds int `json:"expire_seconds,omitempty"` // 二维码的有效时间, 以秒为单位. 最大不超过 604800.
	PermanentQRCode
}

//CreateTemporaryQRCode  创建临时二维码
//  SceneId:       场景值ID, 为32位非0整型
//  ExpireSeconds: 二维码有效时间, 以秒为单位.  最大不超过 604800.
func (c *Qrcode) CreateTemporaryQRCode(SceneID uint32, ExpireSeconds int) (qrcode *TemporaryQRCode, err error) {
	if SceneID == 0 {
		err = errors.New("SceneId should be greater than 0")
		return
	}
	if ExpireSeconds <= 0 {
		err = errors.New("ExpireSeconds should be greater than 0")
		return
	}
	var request struct {
		ExpireSeconds int    `json:"expire_seconds"`
		ActionName    string `json:"action_name"`
		ActionInfo    struct {
			Scene struct {
				SceneID uint32 `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ExpireSeconds = ExpireSeconds
	request.ActionName = "QR_SCENE"
	request.ActionInfo.Scene.SceneID = SceneID

	var result struct {
		mp.Error
		TemporaryQRCode
	}

	response, err := c.HTTPPostJSONWithAccessToken(qrcodeURL, &request)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	result.TemporaryQRCode.SceneID = SceneID
	qrcode = &result.TemporaryQRCode
	return
}

//CreatePermanentQRCode 创建永久二维码
//  SceneId: 场景值ID, 目前参数只支持1--100000
func (c *Qrcode) CreatePermanentQRCode(sceneID uint32) (qrcode *PermanentQRCode, err error) {
	if sceneID == 0 {
		err = errors.New("SceneId should be greater than 0")
		return
	}
	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneID uint32 `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ActionName = "QR_LIMIT_SCENE"
	request.ActionInfo.Scene.SceneID = sceneID

	var result struct {
		mp.Error
		PermanentQRCode
	}

	response, err := c.HTTPPostJSONWithAccessToken(qrcodeURL, &request)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}

	result.PermanentQRCode.SceneID = sceneID
	qrcode = &result.PermanentQRCode
	return
}

//CreatePermanentQRCodeWithSceneString 创建永久二维码
//  SceneString: 场景值ID(字符串形式的ID), 字符串类型, 长度限制为1到64
func (c *Qrcode) CreatePermanentQRCodeWithSceneString(SceneString string) (qrcode *PermanentQRCode, err error) {
	if SceneString == "" {
		err = errors.New("SceneString should not be empty")
		return
	}
	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneString string `json:"scene_str"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ActionName = "QR_LIMIT_STR_SCENE"
	request.ActionInfo.Scene.SceneString = SceneString

	var result struct {
		mp.Error
		PermanentQRCode
	}

	response, err := c.HTTPPostJSONWithAccessToken(qrcodeURL, &request)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}

	result.PermanentQRCode.SceneString = SceneString
	qrcode = &result.PermanentQRCode
	return
}
