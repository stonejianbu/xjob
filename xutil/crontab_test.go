package xutil

import (
	"fmt"
	"testing"
)

func TestStartTimerTask(t *testing.T) {
	count := 0
	StartTimerTask("nice", TS, func() {
		count ++
		fmt.Println("count =", count)
	})
	select {}
}
