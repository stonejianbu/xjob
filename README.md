# xjob

### 介绍
基于`golang+channel`实现的生产消费模型的任务库

### 特性
1. 多任务解耦实现，生产者只管发送消息到任务管道，消息者监听获取任务管道中的消息并进行处理
2. 支持消息处理失败重试，可设置梯级休眠策略
3. 支持消息缓存
4. 支持消息持久化处理，系统异常重启后自动会加载未完成的任务消息到任务管道
5. 支持定时任务
6. 统计和可视化任务状态（TODO）
7. 任务异常告警(TODO)

### 使用
#### 下载安装
```shell
go get github.com/stonejianbu/xjob@v0.1.0
```

#### 开始使用
```go
package main

import (
	"errors"
	"fmt"
	"github.com/stonejianbu/xjob"
	"github.com/stonejianbu/xjob/xmsg"
	"time"
)

func main() {
	// 创建工作对象
	job := xjob.NewDefaultJob()

	// 生产消息
	job.Produce(xmsg.NewMsg("hello", 1022))

	// 消费消息
	job.Consume("hello", func(msg xmsg.Msg) error {
		fmt.Println("-------", msg.Count, msg.Val)
		return errors.New("some go error")
	})

	// 生产消息
	xjob.Produce(xmsg.NewMsg("hello", 1000))

	time.Sleep(time.Second * 50)
}
```

