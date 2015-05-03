package main

import (
    "fmt"
    "git.oschina.net/wida/mycron/src/cron"
    "git.oschina.net/wida/mycron/src/mycron"
    "time"
)

func main() {
    jobs, _ := mycron.GetCronList()
    c := cron.New()
    defer func() {
        c.Stop()
    }()

    //添加jobs
    for i := 0; i < len(jobs); i++ {
        job := jobs[i]
        c.AddFunc(job.Time,
        func() {jobrun(job)},
        int(job.Status), int(job.ID), int64(job.STime), int64(job.ETime))
    }
    //start
    c.Start()

    //监听更新事件
    for {
        select {
        case <-time.After(time.Second):
            jobs, _ := mycron.GetModifyList()
            for i := 0; i < len(jobs); i++ {
                job := jobs[i]
                c.AddFunc(job.Time,
                func() {jobrun(job)},
                int(job.Status), int(job.ID), int64(job.STime), int64(job.ETime))
            }
            mycron.UpdateModifyList()
            continue
        }
    }
}

func jobrun(job mycron.Job){
    defer func() {
        if err := recover(); err != nil {
            mycron.Log(err);
        }
    }()
    job.ChangeRunningStatu(1)
    s, e := mycron.ExecWithTimeout(time.Second*10, job.Cmd)
    if e != nil {
        fmt.Print(e)
    }
    job.ChangeRunningStatu(0)
    fmt.Println(s)
}

/*
func printfEntry(c *cron.Cron) {
    for _, v := range c.Entries() {
        fmt.Println(v.ID, v.Status, v.Start, v.Ending)
    }
}
*/
