package xutil

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

const TM0217 = "0 8 17,2 * *" // 每月2号和17号早上8点执行
const TM15 = "0 0 15 * *"     // 每月15号0点执行
const TD = "@daily"           // 每天0点执行
const TM = "@monthly"         // 每月第一天的0点执行
const TW = "@weekly"          // 每周日开始的那个0点
const TH = "@hourly"          // 每1小时执行
const Tm = "@every 60s"       // 每分钟执行
const TS = "@every 1s"        // 每秒执行

// StartTimerTask 定时任务，添加log和异常处理
func StartTimerTask(taskName, timeStr string, taskFunc func()) {
	c := cron.New(cron.WithLogger(
		cron.VerbosePrintfLogger(log.New(os.Stdout, fmt.Sprintf("%s：", taskName), log.LstdFlags))))
	_, err := c.AddJob(timeStr, cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(&MyJob{Name: taskName, taskFunc: taskFunc}))
	if err != nil {
		log.Println(err)
	}
	c.Start()
	defer c.Stop()
	select {}
}

// MyJob 定时任务工作
type MyJob struct {
	Name     string
	taskFunc func()
}

// Run 开始执行的任务内容
func (that *MyJob) Run() {
	that.taskFunc()
}
