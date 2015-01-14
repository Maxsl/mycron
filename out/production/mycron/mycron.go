package main

import (
	"fmt"
	"github.com/robfig/cron"
	"os/exec"
	"time"
)

const ONE_SECOND = 1*time.Second + 10*time.Millisecond

//exec another process
//if wait d Duration, it will kill the process
//d is <= 0, wait forever
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
	globalchan := make(chan bool)
	c := cron.New()
	defer func() {
		c.Stop()
		close(globalchan)
	}()
	c.AddFunc("* * * * * ?",
		func() {
			fmt.Println("1s")
		})
	c.AddFunc("*/2 * * * * ?",
		func() {
			fmt.Println("2s")
		})
	c.Start()

	time.Sleep(3 * time.Second)
	c.AddFunc("*/3 * * * * ?",
		func() {
			fmt.Println("3s")
		})
	<-globalchan
}
