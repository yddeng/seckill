package seckill

import (
	"github.com/yddeng/seckill/sdk/jd"
	"github.com/yddeng/seckill/util"
	"log"
	"time"
)

func Login() bool {
	// 二维码
	token := ""
	util.LoopFunc(time.Second, func() bool {
		token = jd.QrLoginImage()
		return token != ""
	})

	// 检查二维码扫描状态
	ticket := ""
	util.LoopFunc(time.Second*2, func() bool {
		ticket = jd.QrcodeTicket(token)
		return ticket != ""
	})

	// 检验登陆状态
	if !jd.ValidQRTicket(ticket) || !jd.ValidCookie() {
		log.Println("登录失败")
		return false
	}

	nickName := jd.GetUserInfo()
	log.Println("用户:", nickName, "登陆成功")
	return true
}
