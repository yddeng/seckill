package seckill

import (
	"github.com/BurntSushi/toml"
	"log"
	"time"
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

func GetBuyTimeMs() (int64, string) {
	// 每天12点
	loc, _ := time.LoadLocation("Local")
	now := time.Now()
	nt := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, loc)
	return nt.UnixNano() / 1e6, nt.Format("2006-01-02 15:04:05")
}

func checkTime() {
	buyTimeMs, _ := GetBuyTimeMs()
	if time.Now().UnixNano()/1e6 > buyTimeMs {
		log.Panicln("抢购时间已过")
	}
}
