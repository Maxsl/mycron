package main

import (
	"fmt"
	"cron"
	"mydb"
	"time"
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
				fmt.Println(job.Name, job.Cmd)
			})
	}
	c.Start()
	<-globalchan
}
