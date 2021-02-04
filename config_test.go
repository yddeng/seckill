package seckill

import "testing"

func TestLoadConfig(t *testing.T) {
	_ = LoadConfig("./config.toml")
}
