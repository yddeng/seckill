package util

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os/exec"
	"runtime"
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

func LoopFunc(dur time.Duration, fn func() bool) {
	for !fn() {
		time.Sleep(dur)
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
