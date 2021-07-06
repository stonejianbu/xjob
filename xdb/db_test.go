package xdb

import (
	"fmt"
	"github.com/stonejianbu/xjob/xmsg"
	"testing"
)

func TestBytes(t *testing.T) {
	msg := xmsg.NewMsg("Hello", "hello world")
	fmt.Println("msgID =", msg.ID)
	Set(msg.ID, msg)
	rets := GetAll()
	fmt.Printf("ret == %#v\n", rets)
}
