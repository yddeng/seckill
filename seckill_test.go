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
