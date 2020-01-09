package oauth

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type  Direction struct {
	Ip,RedirectURI, Scope, State string
}

//最快获取微信用户信息的跳转方法
func  (oauth *Oauth) OauthWithCacheInfo(writer http.ResponseWriter, req *http.Request,m Direction, f func(user OauthUser)) {

	agentKey, exist := FilterRedisKeyOfUserAgent(req)
	if !exist {
		_=oauth.Redirect(writer, req, m.RedirectURI, m.Scope,  m.State)
		return
	}

	var wechatUser   OauthUser

	err1 := oauth.Cache.HGet(m.Ip, agentKey, &wechatUser)
	if err1 != nil {
		log.Println(err1)
		_=oauth.Redirect(writer,req,  m.RedirectURI,  m.Scope,  m.State)
		return
	}

	if len(wechatUser.Openid()) == 0 {
		_=oauth.Redirect(writer, req, m.RedirectURI,  m.Scope,  m.State)
		return
	}

		// 取出openid对应的信任度
		weight, err2 := oauth.GetOpenidWeight(wechatUser.Openid())
		if err2 != nil {
		 log.Println(err2)
			_=oauth.Redirect(writer, req,  m.RedirectURI,  m.Scope,  m.State)
			return
		}
		if weight < 50 {
			_=oauth.Redirect(writer, req,  m.RedirectURI,  m.Scope,  m.State)
			return
		}
 

	//触发信任度检查机制
	//1. 更新redis中本人的访问时间
	//2. 5s中之后检查本人的消息回调的记录是否存在，不存在的话，此openid的信任度-20
	//快速登录检查扣分3次，回调检查扣分2次，快速登录扣分1次+回调扣分一次，信任度< 50，将会无法使用快速登录，等到缓存过期又可以重新使用快速登录
	//go ChangeUserOpenidWeight(wx.Openid())

	////直接带上参数重定向到前端页面
	//redirectToFrontWebPage(c, wx, studyCenter)
	f(wechatUser)
}

func (oauth *Oauth) GetOpenidWeight(openid string) (int, error) {

	redisKey := "openidweight:" + openid

	var value string
err := oauth.Cache.GetWithErrorBack(redisKey,value)
	if err != nil {
		return 0, err
	}
	if len(value) == 0 {
		value = "0"
	}

	oldWeight, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return oldWeight, nil

}



func FilterRedisKeyOfUserAgent(req *http.Request) (key string, exist bool) {
	agentStr := req.Header["User-Agent"]
	agentKey := ""
	for _, v := range agentStr {
		str1 := strings.Replace(v, " ", "", -1)
		//str2 := strings.ReplaceAll(str1, "]", "")
		//str3 := strings.ReplaceAll(str2, "{", "")
		//str4 := strings.ReplaceAll(str3, "}", "")
		agentKey += str1
	}
	if len(agentStr) > 0 {
		if len(agentStr) > 1024 {
			//长度太长的redis key，只截取其中一部分
			agentKey = agentKey[0:1024]
		}
		return agentKey, true
	}

	return "", false
}

type OauthUser interface {
	  Openid()string
}