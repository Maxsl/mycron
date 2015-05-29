package mycron

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
    "os/exec"
    "runtime"
    "os"
    "strings"
    "bytes"
    "sync"
)

type Job struct {
	ID              int
	Name, Time, Cmd string
	STime, ETime    int
	Status          uint8
	Running         uint8
	IsModify        uint8
	Process         uint8
	Ip              string
    Singleton       uint8
}

type RunRet struct {
    Pid   int
    Out   string
    Err   error
}

var (
	db  *sql.DB
	err error
)

func init() {
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
                    Mysql_user, Mysql_pwd, Mysql_host, Mysql_prot,Mysql_dbname))
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)
	db.Ping()
}

func GetCronList() (jobss []Job, e error) {
    ut := int64(time.Now().Unix())
	rows, err := db.Query("SELECT id,name,time,cmd,sTime,eTime,status,isrunning,modify,process,ip,singleton FROM cron where status = 1 and sTime < ? and eTime > ?", ut, ut)
	if err != nil {
		panic(err.Error())
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))
	jobs := make([]Job, len(values))
	i := 0
	// Fetch rows
	for rows.Next() {
		err = rows.Scan(&jobs[i].ID, &jobs[i].Name, &jobs[i].Time, &jobs[i].Cmd, &jobs[i].STime, &jobs[i].ETime,
			&jobs[i].Status, &jobs[i].Running,&jobs[i].IsModify,&jobs[i].Process,&jobs[i].Ip,&jobs[i].Singleton)
		if err != nil {
			panic(err.Error())
		}
		i++
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
	return jobs, nil
}

func GetModifyList()(jobss []Job, e error){
    defer func() {
        if err := recover(); err != nil {
            Log(err);
        }
    }()
    ut := int64(time.Now().Unix())
	rows, err := db.Query("SELECT id,name,time,cmd,sTime,eTime,status,isrunning,modify,process,ip ,singleton FROM cron where sTime < ? and eTime > ? and modify = 1", ut, ut)
	if err != nil {
		panic(err.Error())
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))
	jobs := make([]Job, len(values))
	i := 0
	// Fetch rows
	for rows.Next() {
		err = rows.Scan(&jobs[i].ID, &jobs[i].Name, &jobs[i].Time, &jobs[i].Cmd, &jobs[i].STime,
					&jobs[i].ETime, &jobs[i].Status, &jobs[i].Running,&jobs[i].IsModify,&jobs[i].Process,&jobs[i].Ip,&jobs[i].Singleton)
		if err != nil {
			panic(err.Error())
		}
		i++
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
	return jobs, nil
}

func UpdateModifyList() (int64, error){
    ut := int64(time.Now().Unix())
	stmtIns, err := db.Prepare("update cron set modify = 0 where sTime < ? and eTime > ? ")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(ut, ut)
	if err != nil {
		panic(err.Error())
	}
	return result.RowsAffected()
}

func (job Job) ChangeRunningStatus(status int) (int64, error) {
     stmtIns, err := db.Prepare("update cron set isrunning = ? where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(status, job.ID)
	if err != nil {
		panic(err.Error())
	}
	return result.RowsAffected()
}


func (job Job) JobStep(step int,str string,process_id,branch int) (int64,error){
    stmtIns, err := db.Prepare("insert into cron_hist set cId = ?,step = ?,process_id =? ,branch =? ,time = ?,ret=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close()
    result, err := stmtIns.Exec(job.ID, step,process_id,branch,time.Now().Format("2006-01-02 15:04:05"),str)
    if err != nil {
        panic(err.Error())
    }
    return result.RowsAffected()
}

func (job Job) Run(){
    job.ChangeRunningStatus(1)
    job.JobStep(0,"start",0,0)
    if job.Process > 1 { //多进程执行
        wg := new(sync.WaitGroup)
        wg.Add(int(job.Process))
        for i := 0;i< int(job.Process);i++ {
            go func(n int) {
                job.Exec(n+1) //分支id+1
                wg.Done()
            }(i)
        }
        wg.Wait()
    }else{//单进程执行
        job.Exec(0)
    }
    job.ChangeRunningStatus(0)
}

func (job Job) Exec(i int)  {
    var cmd * exec.Cmd
    if runtime.GOOS == "windows"{
        cmd = exec.Command("cmd", "/C", job.Cmd)
    }else {
        shell := os.Getenv("SHELL")
        cmd = exec.Command(shell, "-c", job.Cmd)
    }

    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Start(); err != nil {
        Log(job.ID,cmd.Path,err.Error(),cmd.Process.Pid,i)
        job.JobStep(3,err.Error(),cmd.Process.Pid,i)
        return
    }
    start := "start"
    if i >0 {
        start = "branch start"
    }
    job.JobStep(0,start,cmd.Process.Pid,i)
    done := make(chan error)
    go func() {
        done <- cmd.Wait()
    }()
    select {
    case  err :=<-done:
        if err !=nil{
            Log(job.ID,cmd.Path,err.Error(),cmd.Process.Pid,i)
            job.JobStep(3,err.Error(),cmd.Process.Pid,i)
            return
        }
    }
    fmt.Println(strings.TrimSpace(out.String()))
    job.JobStep(1,strings.TrimSpace(out.String()),cmd.Process.Pid,i)
}