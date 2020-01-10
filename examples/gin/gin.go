package main

import (
	"errors"
	"fmt"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/message"
	"github.com/silenceper/wechat/oauth"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	o := router.Group("/oauth")
	{
		o.GET("/personcenter/begin", PersonCenterOauthBegin)
		o.GET("/personcenter/success", PersonCenterOauthSuccess)
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

//PersonCenterBeginOauth  开始授权登录重定向，获取一次性code
func PersonCenterOauthBegin(c *gin.Context) {
	o := Wc.GetOauth()

	//跳转到授权认证的后接口，先获取微信用户信息，保存微信关注状态，然后绑定微信和游戏账户
	oauthSuccessUrl := getDomainUrl() + "/oauth/personcenter/success"
	m := oauth.Direction{
		Ip:          c.ClientIP(),
		RedirectURI: oauthSuccessUrl,
		Scope:       "snsapi_userinfo",
		State:       "random_string",
	}

	//重定向成功
	o.FastOauthWithCache(c.Writer, c.Request, m, func(user oauth.OauthUser) {
		OperationAfterOauthSuccess(c, user)
	})

}

//PersonCenterRedirectToFrontPage 授权登录拿到了code以后，用code缓存换取微信用户信息，并执行后续操作
func PersonCenterOauthSuccess(c *gin.Context) {

	//通过code换取access_token
	code := c.Query("code")
	if len(code) == 0 {
		log.Println(errors.New("code 为空"))
		c.JSON(http.StatusBadRequest, gin.H{"message": "code 为空"})
		return
	}

	oauth := Wc.GetOauth()

	resToken, err := oauth.GetUserAccessToken(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	//getUserInfo
	userInfo, err := oauth.GetUserInfo(resToken.AccessToken, resToken.OpenID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	OperationAfterOauthSuccess(c, userInfo)
}

func getDomainUrl() string {
	return "http://127.0.0.1"
}

func OperationAfterOauthSuccess(c *gin.Context, user oauth.OauthUser) {
	//保存授权登录获取的用户信息到缓存
	err := saveOauthUserInfoToRedis(c, user)
	if err != nil {
		log.Println(err) //出于高可用，这里并不会return
	}

	//TODO 授权登录成功以后的操作，比如带着用户信息重定向到前端网页

}

func saveOauthUserInfoToRedis(c *gin.Context, user oauth.OauthUser) error {

	agentKey, exist := oauth.FilterRedisKeyOfUserAgent(c.Request)

	if !exist {
		header := fmt.Sprint(c.Request.Header)
		return errors.New("不标准的网络请求 FilterRedisKeyOfUserAgent error: " + header)
	}

	return Wc.GetOauth().Cache.HSetWxUser(c.ClientIP(), agentKey, user)

}
