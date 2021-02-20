package seckill

import (
	"fmt"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	_ = LoadConfig("./config.toml")
}

func TestGetBuyTimeMs(t *testing.T) {
	buyTimeMs, buyTimeStr := GetBuyTimeMs()
	fmt.Println(time.Now().UnixNano()/1e6, buyTimeMs, buyTimeStr)
}
