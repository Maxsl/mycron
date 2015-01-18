package myexec

import (
	"fmt"
	"os/exec"
	"time"
)

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
