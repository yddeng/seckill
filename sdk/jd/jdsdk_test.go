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

func TestGetUserNickname(t *testing.T) {
	nickName := GetUserNickname()
	fmt.Println("nickName", nickName)

	sdk.LoadCookie("./my.cookies")
	nickName = GetUserNickname()
	fmt.Println("nickName", nickName)
}

func TestInitCart(t *testing.T) {
	sdk.LoadCookie("./my.cookies")
	fmt.Println(InitCart("100015185396", "1"))

	// 茅台
	fmt.Println(InitCart("100012043978", "1"))
}

func TestCartIndex(t *testing.T) {
	sdk.LoadCookie("./my.cookies")
	CartIndex()
}

func TestGetSeckillInitInfo(t *testing.T) {
	//initData, err := GetSeckillInitInfo("100012043978", "1")
	initData, err := GetSeckillInitInfo("100009619287", "1")
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

// 预约
func TestGetReserveUrl(t *testing.T) {
	url := GetReserveUrl("100012043978")
	fmt.Println(url)
	if url != "" {
		fmt.Println(RequestReserveUrl(url))
	}
}
