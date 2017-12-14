package gowechat

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/astaxie/beego/cache"
	"github.com/yaotian/gowechat/server"
	"github.com/yaotian/gowechat/server/context"
	"github.com/yaotian/gowechat/util"
)

// Wechat struct
type Wechat struct {
	Context *context.Context
}

// NewWechat init
func NewWechat(cfg context.Config) *Wechat {
	context := new(context.Context)
	initContext(cfg, context)
	return &Wechat{context}
}

func initContext(cfg context.Config, context *context.Context) {
	if cfg.Cache == nil {
		cfg.Cache, _ = cache.NewCache("memory", `{"interval":60}`)
	}
	context.Config = cfg

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

//Mch 商户平台
func (wc *Wechat) Mch() (mch *MchMgr, err error) {
	mch = new(MchMgr)
	mch.Wechat = *wc
	return
}

//Mp 公众平台
func (wc *Wechat) Mp() (mp *MpMgr, err error) {
	err = wc.checkCfgBase()
	if err != nil {
		return
	}
	mp = new(MpMgr)
	mp.Wechat = *wc
	return
}

//checkCfgBase 检查配置基本信息
func (wc *Wechat) checkCfgBase() (err error) {
	return
}

func (wc *Wechat) checkCfgMch() (err error) {
	err = wc.checkCfgBase()
	if err != nil {
		return
	}
	if wc.Context.MchID == "" || wc.Context.MchAPIKey == "" {
		return fmt.Errorf("%s", "配置中没有MchID或者MchAPIKey")
	}
	return
}
