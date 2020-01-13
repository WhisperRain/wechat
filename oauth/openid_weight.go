package oauth

import (
	"github.com/WhisperRain/wechat/util"
	"log"
	"math"
	"time"
)

const OpenIDWeightExpireTime  = 6*30*24*3600*time.Second  //openid的信任度过期时间半年

func (oauth *Oauth) ChangeUserOpenidWeight(loginOpenid string) {
	defer func() {
		nerr := recover()
		if nerr != nil {
			log.Println(nerr.(error))
		}
	}()

	defer func() {
		if e := recover(); e != nil {
			log.Printf("panic error: err=%v", e)
		}
	}()

	redisCache, err := oauth.GetRedisFromCache()
	if err != nil {
		return
	}

	// 5s中之后检查用户的回调的记录是否存在，不存在的话，此openid的信任度-30
	time.Sleep(time.Second * 5)
	redisKey := util.WechatCallBackKeyPrefix + loginOpenid
	var callBackTime int64
	err = redisCache.GetWithErrorBack(redisKey, &callBackTime)

	if err != nil && err.Error() != "redis: nil" {
		//如果数据不存在，那么err==redis: nil
		log.Println(err)
	}

	duration := time.Now().Unix() - callBackTime
    //用户调后台接口和微信回调的时间相差大于10s，则信任度减20
	if math.Abs(float64(duration)) > 10 {
		//openid对应的信任度-20，快速登录检查扣分3次，则信任度< 50，将会刷新缓存
		err := oauth.decreaseOpenidWeight(loginOpenid, 20)
		if err != nil {
			log.Println(err)
		}
	}

}

func (oauth *Oauth) decreaseOpenidWeight(openid string, num int64) error {

	redisCache, err := oauth.GetRedisFromCache()
	if err != nil {
		return err
	}

	redisKey := "openidweight:" + openid

	var oldWeight int
	err = redisCache.GetWithErrorBack(redisKey, &oldWeight)
	if err != nil {
		return err
	}

	if oldWeight < 50 {
		return nil
	}

	err = redisCache.DecrBy(redisKey, num)
	if err != nil {
		return err
	}
	return nil
}

func (oauth *Oauth) getOpenidWeight(openid string) (int, error) {
	redisCache, err := oauth.GetRedisFromCache()
	if err != nil {
		return 0, err
	}

	redisKey := "openidweight:" + openid
	var oldWeight int
	err = redisCache.GetWithErrorBack(redisKey, &oldWeight)
	if err != nil {
		return 0, err
	}

	return oldWeight, nil
}
