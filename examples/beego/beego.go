package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/yaotian/gowechat"
	"github.com/yaotian/gowechat/mp/message"
	"github.com/yaotian/gowechat/mp/oauth"
	"github.com/yaotian/gowechat/wxcontext"
)

var appURL = "http://localhost:8001"
var config = wxcontext.Config{
	AppID:          "your app id",
	AppSecret:      "your app secret",
	Token:          "your token",
	EncodingAESKey: "your encoding aes key",
}

func hello(ctx *context.Context) {
	//配置微信参数
	wc := gowechat.NewWechat(config)
	fmt.Println("wechat:", *wc)
	//微信平台mp
	mp, err := wc.MpMgr()
	if err != nil {
		return
	}

	// 传入request和responseWriter
	msgHandler := mp.GetMsgHandler(ctx.Request, ctx.ResponseWriter)

	//设置接收消息的处理方法
	msgHandler.SetHandleMessageFunc(func(msg message.MixMessage) *message.Reply {
		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{message.MsgTypeText, text}
	})

	//处理消息接收以及回复
	err = msgHandler.Handle()
	if err != nil {
		fmt.Println(err)
	}
}

//wxOAuth 微信公众平台，网页授权
func wxOAuth(ctx *context.Context) {
	//配置微信参数
	wc := gowechat.NewWechat(config)
	//微信公众平台
	mp, err := wc.MpMgr()
	if err != nil {
		return
	}

	oauthHandler := mp.GetPageOAuthHandler(ctx.Request, ctx.ResponseWriter, appURL+"/oauth")

	oauthHandler.SetFuncCheckOpenIDExisting(func(openID string) bool {
		//看自己的系统中是否已经存在此openID的用户
		//如果已经存在， 调用自己的Login 方法，设置cookie等，return true
		//如果还不存在，return false, handler会自动去取用户信息
		return true
	})

	oauthHandler.SetFuncAfterGetUserInfo(func(user oauth.UserInfo) bool {
		//已获得用户信息，这里用信息做注册使用
		//调用自己的Login方法，设置cookie等
		return false
	})

	oauthHandler.Handle()
}

func main() {
	beego.Any("/", hello)
	beego.Any("/oauth", wxOAuth) //需要网页授权的页面url  /oauth?target=url
	beego.Run(":8001")
}
