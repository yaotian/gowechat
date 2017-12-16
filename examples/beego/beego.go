package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/yaotian/gowechat"
	"github.com/yaotian/gowechat/mp/message"
	"github.com/yaotian/gowechat/wxcontext"
)

var appURL = "http://localhost:8001"

func hello(ctx *context.Context) {
	//配置微信参数
	config := wxcontext.Config{
		AppID:          "your app id",
		AppSecret:      "your app secret",
		Token:          "your token",
		EncodingAESKey: "your encoding aes key",
	}
	wc := gowechat.NewWechat(config)

	mp, err := wc.MpMgr()
	if err != nil {
		return
	}

	// 传入request和responseWriter
	msgHandler := mp.GetMsgHandler(ctx.Request, ctx.ResponseWriter)
	fmt.Println("msgHandler:", msgHandler)
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
		return
	}
}

//wxOAuth 微信公众平台，网页授权
func wxOAuth(ctx *context.Context) {
	//配置微信参数
	config := wxcontext.Config{
		AppID:          "your app id",
		AppSecret:      "your app secret",
		Token:          "your token",
		EncodingAESKey: "your encoding aes key",
	}
	wc := gowechat.NewWechat(config)

	mp, err := wc.MpMgr()
	if err != nil {
		return
	}
	oauthHandler := mp.GetPageOAuthHandler(ctx.Request, ctx.ResponseWriter, appURL+"/oauth")
	oauthHandler.Handle()
}

func main() {
	beego.Any("/", hello)
	beego.Any("/oauth", wxOAuth)
	beego.Run(":8001")
}
