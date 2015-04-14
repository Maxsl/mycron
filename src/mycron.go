package main

import (
	"fmt"
	"cron"
	"mydb"
	"time"
	"myexec"
	//"os"
)

const ONE_SECOND = 1*time.Second + 10*time.Millisecond

func main() {
	jobs, _ := mydb.GetCronList()
	globalchan := make(chan bool)
	c := cron.New()
	defer func() {
		c.Stop()
		close(globalchan)
	}()

	for i := 0; i < len(jobs); i++ {
		job := jobs[i]
		c.AddFunc(job.Time,
			func() {
				job.ChangeRunningStatu(1)
				s,e := myexec.ExecWithTimeout(time.Second*10,job.Cmd);
				if e != nil {
					fmt.Print(e)
				}
				job.ChangeRunningStatu(0)
				fmt.Println(s);
			},job.STime,job.ETime)
	}
	c.Start()
	<-globalchan
}
