package seckill

import (
	"fmt"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	_ = LoadConfig("./config.toml")
}

func TestConfig_GetBuyTimeMs(t *testing.T) {
	_ = LoadConfig("./config.toml")
	fmt.Println(time.Now().UnixNano()/1e6, config.GetBuyTimeMs())
}
