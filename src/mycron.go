package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"os/exec"
	"time"
)

type job struct {
	ID              int
	Name, Time, Cmd string
}

func getCronList() (jobss []job, e error) {
	db, err := sql.Open("mysql", "root:@/mycron")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM cron")
	if err != nil {
		panic(err.Error())
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))
	jobs := make([]job, len(values))
	i := 0
	// Fetch rows
	for rows.Next() {
		err = rows.Scan(&jobs[i].ID, &jobs[i].Name, &jobs[i].Time, &jobs[i].Cmd)
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

const ONE_SECOND = 1*time.Second + 10*time.Millisecond

func ExecTimeout(d time.Duration, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if err := cmd.Start(); err != nil {
		return err
	}
	if d <= 0 {
		return cmd.Wait()
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	fmt.Println(cmd.Process.Pid)
	select {
	case <-time.After(d):
		cmd.Process.Kill()
		//wait goroutine return
		//fmt.Println(cmd.Process.Pid)
		return <-done
	case err := <-done:
		return err
	}
}

func main() {
	jobs, _ := getCronList()
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
