package gowechat

import (
	"net/http"
	"sync"

	"github.com/astaxie/beego/cache"
	"github.com/yaotian/gowechat/mch/pay"
	"github.com/yaotian/gowechat/server"
	"github.com/yaotian/gowechat/server/context"
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

	//mch商户平台需要的变量
	//证书
	SslCertFilePath string //证书文件的路径
	SslKeyFilePath  string
	SslCertContent  string //证书的内容
	SslKeyContent   string
	MchID           string
	MchAPIKey       string //商户平台设置的api key
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

	context.MchAPIKey = cfg.MchAPIKey
	context.MchID = cfg.MchID
}

// GetServer 消息管理
func (wc *Wechat) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	wc.Context.Request = req
	wc.Context.Writer = writer
	return server.NewServer(wc.Context)
}

//GetPay get pay
func (wc *Wechat) GetPay() *pay.Pay {
	return pay.NewPay(wc.Context)
}
