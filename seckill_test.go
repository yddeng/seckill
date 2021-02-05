package seckill

import (
	"fmt"
	"testing"
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
