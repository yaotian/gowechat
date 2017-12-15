package main

import (
	"fmt"
	"net/http"

	"github.com/yaotian/gowechat"
	"github.com/yaotian/gowechat/mp/message"
	"github.com/yaotian/gowechat/wxcontext"
)

func hello(rw http.ResponseWriter, req *http.Request) {

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
	msgHandler := mp.GetMsgHandler(req, rw)
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

func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}
