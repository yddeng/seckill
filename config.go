package seckill

import (
	"github.com/BurntSushi/toml"
	"log"
	"time"
)

type Config struct {
	EId string `toml:"eid"`
	Fp  string `toml:"fp"`
	PWD string `toml:"pwd"`

	SkuId   int    `toml:"sku_id"`
	SkuNum  int    `toml:"sku_num"`
	BuyTime string `toml:"buy_time"`
}

func (c *Config) GetBuyTimeMs() int64 {
	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", c.BuyTime, loc)
	return t.UnixNano() / 1e6
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
	if conf.SkuId == 0 || conf.SkuNum == 0 {
		log.Panicln("请填写抢购的商品ID及数量")
	}

	if conf.BuyTime == "" || time.Now().UnixNano()/1e6 > conf.GetBuyTimeMs() {
		log.Panicln("时间格式错误或者时间已过期")
	}

}
