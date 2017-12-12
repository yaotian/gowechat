package gowechat

import (
	"net/http"
	"sync"

	"github.com/astaxie/beego/cache"
	"github.com/yaotian/gowechat/mp/jssdk"
	"github.com/yaotian/gowechat/mp/material"
	"github.com/yaotian/gowechat/mp/menu"
	"github.com/yaotian/gowechat/mp/template"
	"github.com/yaotian/gowechat/mp/user"
	"github.com/yaotian/gowechat/server"
	"github.com/yaotian/gowechat/server/context"
	"github.com/yaotian/gowechat/server/oauth"
	"github.com/yaotian/gowechat/util"
)

// Wechat struct
type Wechat struct {
	Context *context.Context
}

// Config for user
type Config struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
	Cache          cache.Cache

	//证书
	SslCertFilePath string //证书文件的路径
	SslKeyFilePath  string
	SslCertContent  string //证书的内容
	SslKeyContent   string
}

// NewWechat init
func NewWechat(cfg *Config) *Wechat {
	context := new(context.Context)
	initContext(cfg, context)
	return &Wechat{context}
}

func initContext(cfg *Config, context *context.Context) {
	context.AppID = cfg.AppID
	context.AppSecret = cfg.AppSecret
	context.Token = cfg.Token
	context.EncodingAESKey = cfg.EncodingAESKey

	if cfg.Cache == nil {
		cfg.Cache, _ = cache.NewCache("memory", `{"interval":60}`)
	}

	context.SetAccessTokenLock(new(sync.RWMutex))
	context.SetJsAPITicketLock(new(sync.RWMutex))

	//create http client
	if cfg.SslCertFilePath != "" && cfg.SslKeyFilePath != "" {
		if client, err := util.NewTLSHttpClient(cfg.SslCertFilePath, cfg.SslKeyFilePath); err == nil {
			context.SHTTPClient = client
		}
	}

	if cfg.SslCertContent != "" && cfg.SslKeyContent != "" {
		if client, err := util.NewTLSHttpClientFromContent(cfg.SslCertContent, cfg.SslKeyContent); err == nil {
			context.SHTTPClient = client
		}
	}
}

// GetServer 消息管理
func (wc *Wechat) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	wc.Context.Request = req
	wc.Context.Writer = writer
	return server.NewServer(wc.Context)
}

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
