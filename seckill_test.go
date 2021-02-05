package seckill

import (
	"fmt"
	"github.com/yddeng/seckill/util"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	fmt.Println(login())
}

func TestCookieLogin(t *testing.T) {
	fmt.Println(cookieLogin())
}

func TestSeckillSku(t *testing.T) {
	LoadConfig("./config.toml")
	seckillSku("100012043978", "1")
}

func TestDiffTime(t *testing.T) {
	LoadConfig("./config.toml")
	diffTime := getDiffTimeMs()
	buyTimeMs := config.GetBuyTimeMs()
	fmt.Println(config.BuyTime, buyTimeMs, util.GetNowTimeMs(), diffTime)
	time.Sleep(time.Duration(buyTimeMs-util.GetNowTimeMs()-diffTime) * time.Millisecond)
	fmt.Println(time.Now().UnixNano()/1e6, time.Now().String())
}
