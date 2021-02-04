package main

import (
	"github.com/yddeng/seckill"
	"log"
)

func main() {
	seckill.LoadConfig("./config.toml")

	seckill.Seckill()

	log.Println("Seckill Stop")
}
