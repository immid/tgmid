package base

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

type Counter struct {
	index chan int64
}

func (c Counter) Next() int64 {
	next := <-c.index
	c.index <- next + 1
	return next
}

var (
	Version   string
	Verbosity int
)

func LocalExec(command string, arg ...string) (result []string) {
	LogVerbose3("exec: ", command, arg)
	cmd := exec.Command(command, arg...)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		Log(err)
	}
	output := strings.TrimSpace(string(cmdOutput.Bytes()))
	result = strings.Split(output, "\n")
	return result
}

func Log(v ...interface{}) {
	log.Println(v...)
}

func LogVerbose(v ...interface{}) {
	if Verbosity < 1 {
		return
	}
	log.Println(v...)
}

func LogVerbose2(v ...interface{}) {
	if Verbosity < 2 {
		return
	}
	log.Println(v...)
}

func LogVerbose3(v ...interface{}) {
	if Verbosity < 3 {
		return
	}
	log.Println(v...)
}

func NewCounter(index int64) *Counter {
	counter := &Counter{index: make(chan int64, 1)}
	counter.index <- index
	return counter
}
