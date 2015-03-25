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
	//cmd := exec.Command(name, args...)
	shell := os.Getenv("SHELL")
	cmd := exec.Command(shell, "-c", line)
	if err := cmd.Start(); err != nil {
		return "" , err
	}
	if d <= 0 {
		return "",cmd.Wait()
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	//fmt.Println(cmd.Process.Pid)
	select {
	case <-time.After(d):
		cmd.Process.Kill()
		//wait goroutine return
		//fmt.Println(cmd.Process.Pid)
		return "time out ",<-done
	case <-done:
		b, _ := exec.Command(shell, "-c", line).Output()
		return strings.TrimSpace(string(b)),nil
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
