package xjob

import (
	"context"
	"errors"
	"fmt"
	"github.com/stonejianbu/xjob/xdb"
	"github.com/stonejianbu/xjob/xmsg"
	"github.com/stonejianbu/xjob/xutil"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	Lock  sync.Mutex
	RLock sync.RWMutex
)

// Job 定义工作对象结构体
type Job struct {
	Chan      chan xmsg.Msg
	Consumers map[string]func(msg xmsg.Msg) error
}

// NewJob 创建工作对象
func NewJob(ctx context.Context, msgChan chan xmsg.Msg) *Job {
	job := &Job{Chan: msgChan}
	go job.start(ctx)
	return job
}

// NewDefaultJob 创建默认的工作对象
func NewDefaultJob() *Job {
	job := NewJob(context.Background(), NewDefaultChan())
	return job
}

// start 启动任务
func (that *Job) start(ctx context.Context) {
	go func() {
		// 每次程序重启，主动获取硬盘中未处理完成的消息
		for k, v := range xdb.GetAll() {
			xutil.Log().Debugf("[storage] get msg ID=%s from storage", k)
			that.Chan <- v
		}
	}()
	for {
		select {
		case <-ctx.Done():
			xutil.Log().Info("chan task exit")
			return
		case msg := <-that.Chan:
			go that.handler(msg)
		}
	}
}

// handler 任务调度
func (that *Job) handler(msg xmsg.Msg) {
	defer func() {
		if err := recover(); err != nil {
			// 如果消息异常则将从硬盘中删除
			xdb.Del(msg.ID)
			xutil.Log().Debugf("[storage] %s delete", msg.ID)
			xutil.Log().Error(err)
		}
	}()

	if msg.Count == 0 {
		// 如果消息是第一次写入则将其持久化到硬盘中
		xutil.Log().Debugf("[storage] %s persistens storage", msg.ID)
		xdb.Set(msg.ID, msg)
	}
	RLock.Lock()
	method := reflect.ValueOf(that.Consumers[msg.Topic])
	RLock.Unlock()
	xutil.Log().Debugf("[%s] start to consumer val = %#v", msg.Topic, msg.Val)
	values := method.Call([]reflect.Value{reflect.ValueOf(msg)})
	if !values[0].IsNil() {
		xutil.Log().Errorf("[%s] job error：%v", msg.Topic, values[0])
		if msg.Count < msg.MaxCount {
			msg.Count += 1
			that.setTactics(msg)
			that.Chan <- msg
		} else {
			err := errors.New(fmt.Sprintf("[%s] the maximum number of retries has been reached，error：%v", msg.Topic, values[0]))
			xdb.Del(msg.ID)
			xutil.Log().Debugf("[storage] %s delete", msg.ID)
			that.SetAlarm(err)
		}
	} else {
		// 如果消息处理成功则将消息删除
		xdb.Del(msg.ID)
		xutil.Log().Debugf("[storage] %s delete", msg.ID)
	}
}

// setTactics 设置休眠策略，"100,300,500,1000" 单位为毫秒，以英文逗号,间隔
func (that *Job) setTactics(msg xmsg.Msg) {
	times := strings.Split(msg.Tactics, ",")
	if len(times) == 0 {
		return
	}
	// 如果当前重试次数已经大于休眠策略中定义的休眠时间
	if msg.Count > len(times) {
		t, _ := strconv.Atoi(times[len(times)-1])
		time.Sleep(time.Millisecond * time.Duration(t))
		return
	}
	t, _ := strconv.Atoi(times[msg.Count-1])
	time.Sleep(time.Millisecond * time.Duration(t))
}

// SetAlarm 设置告警，邮件、微信或企业微信通知等等
func (that *Job) SetAlarm(err error) {
	xutil.Log().Errorf("异常告警：%s", err.Error())
}

// Consume 注册消费者消费
func (that *Job) Consume(topic string, f func(msg xmsg.Msg) error) {
	if that.Consumers == nil {
		that.Consumers = make(map[string]func(msg xmsg.Msg) error)
	}
	xutil.Log().Debugf("register consumer topic = %s, func = %s", topic, runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
	Lock.Lock()
	that.Consumers[topic] = f
	Lock.Unlock()
}

// Produce 生产者生产消息
func (that *Job) Produce(msg xmsg.Msg)  {
	that.Chan <- msg
}

// Produce 生产者生产消息
func Produce(msg xmsg.Msg) {
	Chan <- msg
}
