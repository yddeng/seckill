package seckill

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	EId string `toml:"eid"`
	Fp  string `toml:"fp"`
	PWD string `toml:"pwd"`
}

var config *Config

func LoadConfig(path string) *Config {
	conf := &Config{}
	_, err := toml.DecodeFile(path, conf)
	if err != nil {
		panic(err)
	}
	config = conf
	log.Println(config)
	return config
}
