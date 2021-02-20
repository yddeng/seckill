package seckill

import "github.com/yddeng/dutil/log"

var logger *log.Logger

func InitLogger(basePath string, fileName string) {
	logger = log.NewLogger(basePath, fileName, 1024*1024*2)
	logger.Debugf("%s logger init", fileName)
}
