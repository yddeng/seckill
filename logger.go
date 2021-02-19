package seckill

import "github.com/yddeng/dutil/log"

var loggerger *log.Logger

func InitLogger(basePath string, fileName string) {
	loggerger = log.NewLogger(basePath, fileName, 1024*1024*2)
	loggerger.Debugf("%s loggerger init", fileName)
}
