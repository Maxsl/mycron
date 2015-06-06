package main
import (
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
        int(job.Status), int(job.Id), int64(job.STime), int64(job.ETime))
    }
    //start
    c.Start()

    //@开启 "立即执行" 监听
    go atonce()

    for {
        //监听更新事件
        select {
            case <-time.After(time.Second):
                jobs, _ := mycron.GetModifyList()
                for _,job:= range jobs{
                    c.AddFunc(job.Time,
                    func() {jobrun(job)},
                    int(job.Status), int(job.Id), int64(job.STime), int64(job.ETime))
                }
                mycron.UpdateModifyList()
                continue
        }
    }
}
//cron执行
func jobrun(job mycron.Job)  {
    defer func() {
        if err := recover(); err != nil {
            mycron.Log(err);
            processSet.Remove(job.Id)
        }
    }()
    if job.Singleton == 1 && processSet.Has(job.Id) { // 如果是单例而且上次还非未退出
        return
    }
    processSet.Add(job.Id)
    job.Run()
    processSet.Remove(job.Id)
}

//立即执行事件处理
func atonce() {
    for {
        //监听更新事件
        select {
        case <-time.After(time.Second):
            jobs, _ := mycron.AtOnce()
            for _, job := range jobs {
                job.Run()
                mycron.UpdateAtOnceList()
                continue
            }
        }
    }
}