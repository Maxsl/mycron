package main
import (
    "fmt"
    "git.oschina.net/wida/mycron/src/cron"
    "git.oschina.net/wida/mycron/src/mycron"
    "time"
)
var(
    processSet =  mycron.NewSet()  //当前正在跑的程序集合
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

func jobrun(job mycron.Job)  {
    defer func() {
        if err := recover(); err != nil {
            mycron.Log(err);
            processSet.Remove(job.ID)
        }
    }()
    if job.Singleton == 1 && processSet.Has(job.ID) { // 如果是单例而且上次还非未退出
        return
    }
    processSet.Add(job.ID)
    job.ChangeRunningStatus(1)
    job.JobStep(0,"start")
    s, e := mycron.ExecWithTimeout(0, job.Cmd)
    job.ChangeRunningStatus(0)
    if e != nil {
        fmt.Print(e)
        processSet.Remove(job.ID)
        job.JobStep(3,e.Error());
        return
    }

    job.JobStep(1,s);
    fmt.Println(s)
    processSet.Remove(job.ID)
}
