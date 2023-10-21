package proxy

import (
	// "fmt"
	"io"
	"log"
	// "net"
	"os"
	"os/exec"
	"sync"

	"github.com/raymond-w-ko/tramp-ssh-redux/pkg/utils"
)

type ServerConnection struct {
	cmd        *exec.Cmd
	stdinPipe  io.WriteCloser
	stdoutPipe io.ReadCloser
}

var serverConns map[string]*ServerConnection = map[string]*ServerConnection{}
var serverConnsLock sync.Mutex

func createServerConn(host string) *ServerConnection {
	var err error
	var conn *ServerConnection = new(ServerConnection)

	conn.cmd = exec.Command("./tramp-ssh-redux-server")

	conn.stdinPipe, err = conn.cmd.StdinPipe()
	if err != nil {
		utils.LogStderr("Could not get stdin pipe for server process " + host)
		log.Fatal(err)
	}

	conn.stdoutPipe, err = conn.cmd.StdoutPipe()
	if err != nil {
		utils.LogStderr("Could not get stdout pipe for server process")
		log.Fatal(err)
	}

	conn.cmd.Stderr = os.Stderr

	err = conn.cmd.Start()
	if err != nil {
		utils.LogStderr("Could not start server process")
		log.Fatal(err)
	}
	utils.LogStderr("Server process started for " + host)

	return conn
}

func getOrCreateServerConn(host string) *ServerConnection {
	serverConnsLock.Lock()
	defer serverConnsLock.Unlock()

	if serverConns[host] == nil {
		var conn *ServerConnection = createServerConn(host)
		serverConns[host] = conn
	}

	return serverConns[host]
}
