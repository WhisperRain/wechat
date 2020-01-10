package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/message"
)

var Wc = wechat.NewWechat(&wechat.Config{
	AppID:          "your app id",
	AppSecret:      "your app secret",
	Token:          "your token",
	EncodingAESKey: "your encoding aes key",
})

func main() {
	router := gin.Default()

	router.Any("/", hello)

 	{
		router.GET("/oauth/personcenter/begin", PersonCenterOauthBegin)
		router.GET(PersoncenterOauthSuccessUrl, PersonCenterOauthSuccess)
	}

	router.Run(":8001")
}

func hello(c *gin.Context) {

	// 传入request和responseWriter
	server := Wc.GetServer(c.Request, c.Writer)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {

		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}
