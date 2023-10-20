package proxy

import (
	// "bytes"
	"encoding/binary"
	"fmt"
	"log"
	"io"
	// "net"
	"os"
	"os/exec"

	"github.com/raymond-w-ko/tramp-ssh-redux/pkg/utils"
)

var serverCmd *exec.Cmd
var serverStdinPipe io.WriteCloser
var serverStdoutPipe io.ReadCloser

func SetupServerProcess() {
	var err error
	serverCmd := exec.Command("./tramp-ssh-redux-server")

	serverStdinPipe, err = serverCmd.StdinPipe()
	if err != nil {
		utils.LogStderr("Could not get stdin pipe for server process")
		log.Fatal(err)
	}

	serverStdoutPipe, err = serverCmd.StdoutPipe()
	if err != nil {
		utils.LogStderr("Could not get stdout pipe for server process")
		log.Fatal(err)
	}

	serverCmd.Stderr = os.Stderr

	err = serverCmd.Start()
	if err != nil {
		utils.LogStderr("Could not start server process")
		log.Fatal(err)
	}
	fmt.Println("Server process started")
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestServerEcho() {
	fmt.Println("Testing server echo")
	id := uint64(0)
	binary.Write(serverStdinPipe, binary.BigEndian, id)
	s := "Hello, world!"
	data := []byte(s)
	n := uint64(len(data))
	binary.Write(serverStdinPipe, binary.BigEndian, n)
	serverStdinPipe.Write(data)
}
