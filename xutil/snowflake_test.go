package xutil

import (
	"fmt"
	"testing"
)

func TestWorker_GetId(t *testing.T) {
	ret := NewWorker(1).GetId()
	fmt.Printf("ret == %#v\n", ret)
}
