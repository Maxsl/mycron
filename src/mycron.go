package main

import (
	"fmt"
	"cron"
	"mydb"
	"time"
	"myexec"
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
				//fmt.Println(job.Name, job.Cmd)
				e := myexec.ExecTimeout(time.Second*10, "/bin/sh", "-c", `ps -ef | grep -v "grep" | grep "php" >> /home/wida/test.txt`)
				if e != nil {
					fmt.Print(e)
				}
			})
	}
	c.Start()
	<-globalchan
}
