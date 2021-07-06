package xjob

import (
	"github.com/stonejianbu/xjob/xmsg"
	"sync"
)

// Chan 任务管道
var Chan chan xmsg.Msg
var DefaultLen = 2000

// NewChan 创建任务管道
func NewChan(chanLength int) chan xmsg.Msg {
	var once sync.Once
	once.Do(func() {
		Chan = make(chan xmsg.Msg, chanLength)
	})
	return Chan
}

// NewDefaultChan 创建默认的任务管道
func NewDefaultChan() chan xmsg.Msg {
	return NewChan(DefaultLen)
}
