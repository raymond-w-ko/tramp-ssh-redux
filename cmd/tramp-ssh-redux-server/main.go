package main

import (
	// "os"
	// "fmt"
	"github.com/raymond-w-ko/tramp-ssh-redux/pkg/utils"
)

type Tasklet struct {
	id uint64
	in chan []byte
	out chan []byte
}

var idToTasklet = map[uint64]*Tasklet{}

func main() {
	go handleStdout()
	handleStdin()
}

func handleStdin() {
	taskletId := utils.ReadStdinUint64()
	msgSize := utils.ReadStdinUint64()
	msg := utils.ReadStdinN(msgSize)

	if (taskletId == 0) {
		s := string(msg)
		utils.LogStderr("server: from stdin:")
		utils.LogStderr(s)
	}
}

func handleStdout() {
}
