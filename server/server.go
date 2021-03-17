package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"syscall"
	cron2 "task-schedule-center/server/cron"
	"time"
)

func main() {
	c := cron.New()

	cron2.InitTask(c)

	c.Start()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range ch {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
				//退出
				c.Stop()
				fmt.Println("clean complete...")
				os.Exit(0)
			case syscall.SIGUSR1:
				fmt.Println("usr1", s)
			case syscall.SIGUSR2:
				fmt.Println("usr2", s)
			default:
				fmt.Println("other signal", s)
			}
		}
	}()
	fmt.Println("server started...")
	for {
		time.Sleep(time.Second)
	}
}
