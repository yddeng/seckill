package main

import "github.com/yddeng/seckill"

func main() {
	seckill.LoadConfig("./config.toml")

	seckill.Login()
}
