package util

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"
	"time"
)

func OpenImage(file string) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("start", file)
		_ = cmd.Start()
	} else {
		if runtime.GOOS == "linux" {
			cmd := exec.Command("eog", file)
			_ = cmd.Start()
		} else {
			cmd := exec.Command("open", file)
			_ = cmd.Start()
		}
	}
}

func LoopFunc(fn func() bool, dur ...time.Duration) {
	sleepDur := time.Duration(0)
	if len(dur) > 0 {
		sleepDur = dur[0]
	}
	for !fn() {
		if sleepDur != 0 {
			time.Sleep(sleepDur)
		}
	}
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Go(num int, fn func(int)) {
	for i := 1; i <= num; i++ {
		id := i
		go fn(id)
	}
}

func WaitGoLoop(goNum int, endTime time.Time, fn func(chan interface{}) bool) interface{} {
	out := make(chan interface{}, goNum)
	ok := int32(0)
	Go(goNum, func(id int) {
		sleepTime := time.Duration(rand.Int()%10+id*5) * time.Millisecond
		LoopFunc(func() bool {
			if fn(out) {
				atomic.StoreInt32(&ok, 1)
				return true
			} else if atomic.LoadInt32(&ok) == 1 {
				// 已经有其他线程完成
				return true
			} else if time.Now().After(endTime) {
				// 超时
				out <- nil
				return true
			}

			return false
		}, sleepTime)
	})

	return <-out
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func GetNowTimeMs() int64 {
	return time.Now().UnixNano() / 1e6
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
