package util

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestWaitGoLoop(t *testing.T) {
	f := func() int64 {
		if rand.Int()%5 == 0 {
			return time.Now().Unix()
		}
		return 0
	}

	nowUnix := WaitGoLoop(2, time.Now().Add(time.Second), func(i chan interface{}) bool {
		t := f()
		if t != 0 {
			i <- t
			return true
		}
		return false
	})

	fmt.Println(nowUnix)

}
