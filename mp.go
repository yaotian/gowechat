// @description  微信公共平台的接口

package gowechat

import (
	"github.com/yaotian/gowechat/mp/jssdk"
	"github.com/yaotian/gowechat/mp/material"
	"github.com/yaotian/gowechat/mp/menu"
	"github.com/yaotian/gowechat/mp/template"
	"github.com/yaotian/gowechat/mp/user"
	"github.com/yaotian/gowechat/server/oauth"
)

//GetAccessToken 获取access_token
func (wc *Wechat) GetAccessToken() (string, error) {
	return wc.Context.GetAccessToken()
}

// GetOauth oauth2网页授权
func (wc *Wechat) GetOauth() *oauth.Oauth {
	return oauth.NewOauth(wc.Context)
}

// GetMaterial 素材管理
func (wc *Wechat) GetMaterial() *material.Material {
	return material.NewMaterial(wc.Context)
}

// GetJs js-sdk配置
func (wc *Wechat) GetJs() *jssdk.Js {
	return jssdk.NewJs(wc.Context)
}

// GetMenu 菜单管理接口
func (wc *Wechat) GetMenu() *menu.Menu {
	return menu.NewMenu(wc.Context)
}

// GetUser 用户管理接口
func (wc *Wechat) GetUser() *user.User {
	return user.NewUser(wc.Context)
}

// GetTemplate 模板消息接口
func (wc *Wechat) GetTemplate() *template.Template {
	return template.NewTemplate(wc.Context)
}
