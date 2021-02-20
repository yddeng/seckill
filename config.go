package seckill

import (
	"github.com/BurntSushi/toml"
	"log"
)

const (
	CookieFilename  = "./my.cookies"
	QrImageFilename = "./qr_code.png"
)

type Config struct {
	EId string `toml:"eid"`
	Fp  string `toml:"fp"`
	PWD string `toml:"pwd"`

	SkuId  string `toml:"sku_id"`
	SkuNum string `toml:"sku_num"`
}

var config *Config

func LoadConfig(path string) *Config {
	conf := &Config{}
	_, err := toml.DecodeFile(path, conf)
	if err != nil {
		panic(err)
	}
	checkConfig(conf)

	config = conf
	log.Println(config)
	return config
}

func checkConfig(conf *Config) {
	if conf.EId == "" || conf.Fp == "" {
		log.Panicln("请填写eid，fp")
	}
	if conf.SkuId == "" || conf.SkuNum == "" {
		log.Panicln("请填写抢购的商品ID及数量")
	}

}
