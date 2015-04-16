package main

import (
	"fmt"
	"git.oschina.net/wida/mycron/src/cron"
	"git.oschina.net/wida/mycron/src/mydb"
	"git.oschina.net/wida/mycron/src/myexec"
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
				job.ChangeRunningStatu(1)
				s, e := myexec.ExecWithTimeout(time.Second*10, job.Cmd)
				if e != nil {
					fmt.Print(e)
				}
				job.ChangeRunningStatu(0)
				fmt.Println(s)
			},int(job.Status),int(job.ID), int64(job.STime), int64(job.ETime))
	}

/*	for _,v := range c.Entries(){
		fmt.Println(v.ID,v.Status,v.Start,v.Ending)
	}*/
	c.Start()

	for {
		now := time.Now().Local()
		select {
		case now = <-time.After(time.Second * 60):
			jobs, _ := myexec.GetModifyList()
			for i := 0; i < len(jobs); i++ {
				job := jobs[i]
				c.AddFunc(job.Time,
					func() {
						job.ChangeRunningStatu(1)
						s, e := myexec.ExecWithTimeout(time.Second*10, job.Cmd)
						if e != nil {
							fmt.Print(e)
						}
						job.ChangeRunningStatu(0)
						fmt.Println(s)
					},int(job.Status),int(job.ID), int64(job.STime), int64(job.ETime))
			}
			myexec.UpdateModifyList()
			continue
		}
	}

	<-globalchan
}
