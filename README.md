# WeChat SDK for Go
[!原始仓库及基本说明](https://github.com/silenceper/wechat)

阅读本文之前,相信您应该看过了读过了原始仓库silenceper/wechat相关说明 。

## 本仓库新实现的一些功能
1.发送客服消息
2.业内第一个实现,无需输入手机号,即可在公众号内绑定手机号. 前提是,用户是从第三方app引导到微信的, 目的也是绑定公众号和第三方app中的手机号.


```go
//测试发送客服消息
func TestCustomerMessage(t *testing.T) {

	//存在一个48小时互动的要求
	go func() {
		mes := message.NewCustomerTextMessage(
			"o2Kps1A63BRkfYf6rPkJXkDd0_i5",
			"欢迎我们的公众号~\n",
		)

		err1 := service.MessageManager.Send(mes)
		if err1 != nil {
			t.Fatal(err1)
		}
	}()

	go func() {
		mes := message.NewCustomerImgMessage(
			"o2Kps1A63BRkfYf6rPkJXkDd0_i5",
			"abcEtCVN5n_6yPgDxNXBQP2mEmBdRI4aUWYNaJPXkDw",
		)

		err2 := service.MessageManager.Send(mes)
		if err2 != nil {
			t.Fatal(err2)
		}
	}()

}


```
完整代码：[message/customer_message.go](./message/customer_message.go)
 

**公众号自动绑定手机号的流程**

所谓自动绑定手机号,其实就是在关注公众号的过程中绑定手机号. 当然也就有一个引导关注公众号的流程. 当一个第三方app引导用户关注公众号的时候,也就自动绑定这个app中用户注册的手机号.

引导关注及绑定的流程
1.获取手机号
    app内,用户登录账户以后,获取用户的注册的手机号. (也可以是userID)
    
2.带着手机号跳转到微信app
    app内部引导用户授权,使用微信开放平台的相关配置,给用户发送微信一次性订阅消息, 用户根据提示跳转到微信app内.这条一次性订阅消息中,有消息详情页的url.我们将手机号拼接到该url中.
    
3.带着手机号跳转到下一个url
    在微信app内点击开放平台的一次性订阅消息, 跳转到消息详情页接口. 我们再引导用户再发送微信公众号的一次性订阅消息. 这是第二条一次性订阅消息.这次使用公众号的配置.我们在这个接口url里面可以获取手机号, 然后我们把手机号传递给下一个url,也就是传递给公众号一次性订阅消息重定向的页面.

4.在url中获取手机号, 公众号的openid, 并绑定
    第二条一次性订阅消息,用户确认接受以后.将会重定向到目标url, 微信会在query中加上参数openid. 由于url中本来就有手机号,便可以在这个接口中绑定手机号和openid.
    
5.跳转到公众号历史消息列表,完成关注
    绑定了openid和手机号以后,在这个接口里,重定向到一个永久微信图文素材,这个图文素材中只需要一个按钮,写着关注公众号.点击按钮跳转到公众号的历史消息列表.最后点击按钮关注,完成关注公众号.

总结
这样引导关注并绑定, 会把关注的流程稍微弄长一点,但是会完全缩减和省略绑定的流程.


## License

Apache License, Version 2.0
