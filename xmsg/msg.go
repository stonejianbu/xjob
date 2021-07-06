package xmsg

import (
	"fmt"
	"github.com/stonejianbu/xjob/xutil"
)

// DefaultMode 失败重试的休眠策略
var DefaultMode = "300,500,1000,3000,5000"
var MinMode = "300"
var MidMode = "500,1000,3000,5000,10000,30000"
var MaxMode = "1000,3000,5000,10000,30000,60000,300000,600000,3600000"

// Msg 通知消息对应的结构体
type Msg struct {
	ID       string      // 消息ID
	Val      interface{} // 消息值
	Topic    string      // 消息主题
	Count    int         // 计数
	MaxCount int         // 最大重试次数
	Tactics  string      // 休眠策略，如"100,300,500" 第一次是100ms,第二次是300ms，第三次是500ms
}

// NewMsg 创建管道消息实例
func NewMsg(topic string, val interface{}) Msg {
	return Msg{
		ID:       fmt.Sprintf("%d", xutil.NewWorker(1).GetId()),
		Val:      val,
		Topic:    topic,
		Count:    0,
		MaxCount: 20,
		Tactics:  MinMode,
	}
}

// NewMsgWithTactics 创建消息实例并设置休眠策略
func NewMsgWithTactics(topic string, val interface{}, tactics string) Msg {
	return Msg{
		ID:       fmt.Sprintf("%d", xutil.NewWorker(1).GetId()),
		Val:      val,
		Topic:    topic,
		Count:    0,
		MaxCount: 20,
		Tactics:  tactics,
	}
}

// SetTactics 设置消息梯级休眠策略，如"100,300,500" 第一次是100ms,第二次是300ms，第三次是500ms
func (that Msg) SetTactics(tactics string) Msg {
	that.Tactics = tactics
	return that
}

// SetMaxCount 设置消息的最大重试次数
func (that Msg) SetMaxCount(max int) Msg {
	if max < 1 {
		max = 1
	}
	that.MaxCount = max
	return that
}