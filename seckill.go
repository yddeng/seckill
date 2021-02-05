package seckill

import (
	"fmt"
	"github.com/yddeng/seckill/sdk"
	"github.com/yddeng/seckill/sdk/jd"
	"github.com/yddeng/seckill/util"
	"log"
	"runtime"
	"time"
)

func cookieLogin() bool {
	if util.Exists(CookieFilename) {
		sdk.LoadCookie(CookieFilename)
		if jd.ValidCookie() {
			nickName := jd.GetUserNickname()
			log.Println(nickName, "本地cookie 登录成功")
			return true
		}
		log.Println("本地cookie 过期")
		return false
	}
	return false
}

func login() bool {
	// 二维码
	token := ""
	util.LoopFunc(func() bool {
		token = jd.QrLoginImage(QrImageFilename)
		return token != ""
	}, time.Second)

	// 检查二维码扫描状态
	ticket := ""
	util.LoopFunc(func() bool {
		ticket = jd.QrcodeTicket(token)
		return ticket != ""
	}, time.Second*2)

	// 检验登陆状态
	if !jd.ValidQRTicket(ticket) || !jd.ValidCookie() {
		log.Println("登录失败")
		return false
	}

	// 保存cookie
	sdk.SaveCookie(CookieFilename)

	nickName := jd.GetUserNickname()
	log.Println("用户:", nickName, "登陆成功")
	return true
}

func seckillSku(skuId, skuNum string) {
	goNum := runtime.NumCPU()
	// 结束时间
	endTime := time.Now().Add(time.Second * 10)

	log.Println(" -- SeckillSku Step 1 -- ")
	killUrl := util.WaitGoLoop(goNum, endTime, func(i chan interface{}) bool {
		killUrl := jd.GetKillUrl(skuId)
		if killUrl != "" {
			i <- killUrl
			return true
		}
		return false
	})
	if killUrl == nil {
		log.Println("超时结束")
		return
	}

	log.Println(" -- SeckillSku Step 2 -- ")
	killUrlReq := util.WaitGoLoop(goNum, endTime, func(i chan interface{}) bool {
		ok := jd.RequestKillUrl(skuId, killUrl.(string))
		if ok {
			i <- true
			return true
		}
		return false
	})
	if killUrlReq == nil {
		log.Println("超时结束")
		return
	}

	log.Println(" -- SeckillSku Step 3 -- ")
	seckillPageReq := util.WaitGoLoop(goNum, endTime, func(i chan interface{}) bool {
		ok := jd.SeckillPage(skuId, killUrl.(string))
		if ok {
			i <- true
			return true
		}
		return false
	})
	if seckillPageReq == nil {
		log.Println("超时结束")
		return
	}

	log.Println(" -- SeckillSku Step 4 -- ")
	initData := util.WaitGoLoop(goNum, endTime, func(i chan interface{}) bool {
		initData, err := jd.GetSeckillInitInfo(skuId, skuNum)
		if err == nil {
			i <- initData
			return true
		}
		return false
	})
	if initData == nil {
		log.Println("超时结束")
		return
	}

	log.Println(" -- SeckillSku Step 5 -- ")
	util.WaitGoLoop(goNum, endTime, func(i chan interface{}) bool {
		jd.SubmitSeckillOrder(config.EId, config.Fp, skuId, skuNum, config.PWD, initData.(*jd.InitData))
		return false
	})
}

func getDiffTimeMs() int64 {
	// 请求开始时间
	nowTimeMs := time.Now().UnixNano() / 1e6
	// 服务器时间
	serverTimeMs := jd.GetServerTime()
	// 请求到服务器花费的时间
	return serverTimeMs - nowTimeMs
}

func Seckill() {

	if !cookieLogin() {
		login()
	}

	log.Println(fmt.Sprintf("等待到达抢购时间:%s", config.BuyTime))

	buyTimeMs := config.GetBuyTimeMs()
	if buyTimeMs-util.GetNowTimeMs() > 60*1000 {
		// 提前60s唤醒
		time.Sleep(time.Millisecond * time.Duration(buyTimeMs-util.GetNowTimeMs()-60*1000))
		// 检查过期
		if !jd.ValidCookie() {
			log.Println("cookie过期, 请重新登陆！")
			return
		}
	}

	diffTime := getDiffTimeMs()
	log.Println(fmt.Sprintf("等待到达抢购时间:%s，检测本地时间与京东服务器时间误差为【%d】毫秒", config.BuyTime, diffTime))
	// 提前500毫秒执行
	time.Sleep(time.Duration(buyTimeMs-diffTime-util.GetNowTimeMs()-500) * time.Millisecond)

	log.Println("时间到达，开始执行……")
	seckillSku(config.SkuId, config.SkuNum)
}
