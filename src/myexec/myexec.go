package myexec

import (
	//"fmt"
	"os/exec"
	"time"
	"os"
	"strings"
	"errors"
)

func ExecWithTimeout(d time.Duration, line string)(string, error) {
	shell := os.Getenv("SHELL")
	cmd := exec.Command(shell, "-c", line)
	if err := cmd.Start(); err != nil {
		return "" , err
	}
	if d <= 0 {
		cmd.Wait()
		b, err := exec.Command(shell, "-c", line).Output()
		return strings.TrimSpace(string(b)),err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	//fmt.Println(cmd.Process.Pid)
	select {
	case <-time.After(d):
		cmd.Process.Kill()
		return "",errors.New("time out")
	case  err :=<-done:
		if err !=nil{
			return "",err
		}
		b, err := exec.Command(shell, "-c", line).Output()
		return strings.TrimSpace(string(b)),err
	}
}

func ShellRun(line string) (string, error) {
	shell := os.Getenv("SHELL")
	b, err := exec.Command(shell, "-c", line).Output()
	if err != nil {
		return "", errors.New(err.Error() + ":" + string(b))
	}
	return strings.TrimSpace(string(b)), nil
}
