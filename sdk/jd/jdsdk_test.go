package jd

import (
	"fmt"
	"testing"
	"time"
)

func TestQrLoginImage(t *testing.T) {
	token := QrLoginImage()
	fmt.Println("token wlfstkSmdl", token)
}

func TestQrcodeTicket(t *testing.T) {
	token := QrLoginImage()
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
