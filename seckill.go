package seckill

import (
	"fmt"
	"github.com/yddeng/dutil"
	"github.com/yddeng/seckill/sdk"
	"github.com/yddeng/seckill/sdk/jd"
	"github.com/yddeng/seckill/util"
	"os"
	"runtime"
	"time"
)

func cookieLogin() bool {
	if util.Exists(CookieFilename) {
		logger.Infoln("验证本地cookie...")
		sdk.LoadCookie(CookieFilename)
		if jd.ValidCookie() {
			nickName := jd.GetUserNickname()
			logger.Infoln(nickName, "本地cookie 登录成功")
			return true
		}
		logger.Infoln("本地cookie 过期")
		return false
	}
	return false
}

func login() bool {
	logger.Infoln("用户登陆流程...")
	// 二维码
	token := ""
	util.LoopFunc(func() bool {
		token = jd.QrLoginImage(QrImageFilename)
		return token != ""
	}, time.Second)

	logger.Infoln("二维码获取成功，请打开京东APP扫描")
	util.OpenImage(QrImageFilename)

	// 检查二维码扫描状态
	ticket := ""
	util.LoopFunc(func() bool {
		ticket = jd.QrcodeTicket(token)
		return ticket != ""
	}, time.Second*2)

	logger.Infoln("已完成手机客户端确认")

	// 检验登陆状态
	if !jd.ValidQRTicket(ticket) || !jd.ValidCookie() {
		logger.Infoln("登录失败")
		return false
	}

	// 保存cookie
	sdk.SaveCookie(CookieFilename)

	nickName := jd.GetUserNickname()
	logger.Infoln("用户:", nickName, "登陆成功")
	return true
}

func seckillSku(skuId, skuNum string) {
	logger.Infoln("执行秒杀抢购流程...")
	goNum := runtime.NumCPU()
	// 结束时间
	endTime := time.Now().Add(time.Second * 60)

	exitFunc := func(i interface{}) {
		if i == nil {
			logger.Infoln("任务超时，程序结束")
			os.Exit(0)
		}
	}

	logger.Infoln("Step1 -- 获取秒杀链接... ")
	killUrl := util.WaitGoLoop(goNum, endTime, func(i chan interface{}) bool {
		killUrl := jd.GetKillUrl(skuId)
		if killUrl != "" {
			i <- killUrl
			util.Go(1, func(i int) {
				defer func() {
					err := dutil.Recover()
					if err != nil {
						logger.Errorln(err)
					}
				}()
				jd.RequestKillUrl(skuId, killUrl)
				jd.RequestSeckillPage(skuId, skuNum)
			})
			return true
		}
		return false
	})
	exitFunc(killUrl)
	logger.Infoln("Step1 --", killUrl)

	logger.Infoln("Step2 -- 获取秒杀商品初始化信息... ")
	initData := util.WaitGoLoop(goNum, endTime, func(i chan interface{}) bool {
		initData, err := jd.GetSeckillInitInfo(skuId, skuNum)
		if err == nil {
			i <- initData
			return true
		}
		return false
	})
	exitFunc(initData)
	logger.Infoln("Step2 --", initData)

	logger.Infoln("Step3 -- 提交秒杀商品订单... ")
	util.WaitGoLoop(goNum, endTime, func(i chan interface{}) bool {
		jd.SubmitSeckillOrder(config.EId, config.Fp, skuId, skuNum, config.PWD, initData.(*jd.InitData))
		return false
	})
	logger.Infoln("Step3 -- OK")

	logger.Infoln("All Steps OK")
}

func getDiffTimeMs() int64 {
	// 请求开始时间
	nowTimeMs := time.Now().UnixNano() / 1e6
	// 服务器时间
	serverTimeMs := jd.GetServerTime()
	// 请求到服务器花费的时间
	if serverTimeMs != 0 {
		return serverTimeMs - nowTimeMs
	}
	return 0
}

func Seckill() {

	// 时间判断
	buyTimeMs, buyTimeStr := util.GetTodayTimeMs(12, 0, 0)
	if util.GetNowTimeMs() > buyTimeMs {
		logger.Infoln("抢购时间已过，明天再来")
		return
	}

	if !cookieLogin() && !login() {
		logger.Infoln("用户登陆失败！！")
		return
	}

	// 预约
	reserveTimeMs, reserveTimeStr := util.GetTodayTimeMs(10, 0, 0)
	if util.GetNowTimeMs() < reserveTimeMs {
		// 预约时间 1分钟后
		logger.Infoln(fmt.Sprintf("等待到达预约时间:%s，将在预约时间后60s唤醒", reserveTimeStr))
		time.Sleep(time.Duration(reserveTimeMs-util.GetNowTimeMs()+60*1000) * time.Millisecond)
	}

	util.Go(1, func(i int) {
		util.LoopFunc(func() (ok bool) {
			defer func() {
				err := dutil.Recover()
				if err != nil {
					logger.Errorln(err)
					ok = true
				}
			}()
			reserveUrl := jd.GetReserveUrl(config.SkuId)
			ok = reserveUrl != ""
			if ok {
				jd.RequestReserveUrl(reserveUrl)
			}
			return
		}, time.Second)
	})

	if buyTimeMs-util.GetNowTimeMs() > 60*1000 {
		// 提前60s唤醒
		logger.Infoln(fmt.Sprintf("等待到达抢购时间:%s，将在开始前60s唤醒", buyTimeStr))
		time.Sleep(time.Millisecond * time.Duration(buyTimeMs-util.GetNowTimeMs()-60*1000))
		// 检查过期
		if !jd.ValidCookie() {
			logger.Infoln("cookie过期, 请重新登陆！")
			return
		}
	}

	diffTime := getDiffTimeMs()
	logger.Infoln(fmt.Sprintf("等待到达抢购时间:%s，检测本地时间与京东服务器时间误差为【%d】毫秒", buyTimeStr, diffTime))
	// 提前500毫秒执行
	time.Sleep(time.Duration(buyTimeMs-diffTime-util.GetNowTimeMs()-500) * time.Millisecond)

	seckillSku(config.SkuId, config.SkuNum)
}
