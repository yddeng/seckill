package jd

import (
	"fmt"
	"github.com/yddeng/seckill/sdk"
	"testing"
	"time"
)

func TestQrLoginImage(t *testing.T) {
	token := QrLoginImage("./qr_code.png")
	fmt.Println("token wlfstkSmdl", token)
}

func TestQrcodeTicket(t *testing.T) {
	token := QrLoginImage("./qr_code.png")
	fmt.Println("token wlfstkSmdl", token)

	var ticket string
	for ticket == "" {
		time.Sleep(time.Second * 2)
		ticket = QrcodeTicket(token)
	}

	fmt.Println("ticket", ticket)
}

func TestGetServerTime(t *testing.T) {
	nowTimeMs := time.Now().UnixNano() / 1e6
	serverTimeMs := GetServerTime()
	diffTime := serverTimeMs - nowTimeMs

	fmt.Println(nowTimeMs, serverTimeMs, diffTime)
}

func TestValidCookie(t *testing.T) {
	sdk.LoadCookie("./my.cookies")
	fmt.Println(ValidCookie())
}

func TestGetKillUrl(t *testing.T) {
	killUrl := GetKillUrl("100012043978")
	fmt.Println(killUrl)
}

func TestGetUserInfo(t *testing.T) {
	sdk.LoadCookie("./my.cookies")
	nickName := GetUserInfo()
	fmt.Println("nickName", nickName)
}

func TestGetSeckillInitInfo(t *testing.T) {
	initData, err := GetSeckillInitInfo("100012043978", "1")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(initData)

}

func TestSubmitSeckillOrder(t *testing.T) {
	initData, err := GetSeckillInitInfo("100012043978", "1")
	if err != nil {
		fmt.Println(err)
		return
	}
	ok := SubmitSeckillOrder("eid", "fp", "skuId", "skuNum", "pwd", initData)
	fmt.Println("submit", ok)
}
