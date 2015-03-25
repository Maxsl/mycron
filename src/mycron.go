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
				//fmt.Println(job.Name, job.Cmd)
			//	e := myexec.ExecTimeout(time.Second*10, os.Getenv("SHELL"), "-c", `ps -ef | grep -v "grep" | grep "php" >> /home/wida/test.txt`)
				//e := myexec.ExecTimeout(time.Second*10, os.Getenv("SHELL"), "-c", "/home/wida/sh.sh")
				s,e := myexec.ExecWithTimeout(time.Second*10,"/home/wida/sh.sh");
				if e != nil {
					fmt.Print(e)
				}
				fmt.Println(s);
			})
	}

//	for _,e:=  range c.Entries(){
//		fmt.Printf("#%v",e.Next)
//		fmt.Printf("#%v",e.Prev)
//		fmt.Printf("#%v",e.Schedule)
//		//fmt.Printf()
//	}

	c.Start()
	<-globalchan
}
