package xjob

import (
	"errors"
	"fmt"
	"github.com/stonejianbu/xjob/xmsg"
	"testing"
	"time"
)

func TestJob(t *testing.T) {
	job := NewDefaultJob()
	job.Consume("Hello", func(msg xmsg.Msg) error {
		fmt.Println("===============nice")
		return errors.New("no")
	})

	Produce(xmsg.NewMsg("Hello", 1526))
	time.Sleep(time.Second*20)
}