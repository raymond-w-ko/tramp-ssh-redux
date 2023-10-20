package proxy

import (
	"os"
	"log"
	"net"
	"github.com/raymond-w-ko/tramp-ssh-redux/pkg/utils"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func SetupUnixSocket() {
	os.Remove(utils.SocketFile)

	addr, err := net.ResolveUnixAddr("unix", utils.SocketFile)
	if err != nil {
		utils.LogStderr("Could open addr for UNIX socket " + utils.SocketFile)
		log.Fatal(err)
	}

	listener, err := net.ListenUnix("unix", addr)
	if err != nil {
		utils.LogStderr("Could not listen on UNIX socket " + utils.SocketFile)
		log.Fatal(err)
	}
	defer os.Remove(utils.SocketFile)
	defer listener.Close()
	utils.LogStderr("Listening on UNIX socket " + utils.SocketFile)

	for {
		conn, err := listener.AcceptUnix()
		if err != nil {
			utils.LogStderr("Could not accept connection on UNIX socket " + utils.SocketFile)
			// should this be fatal?
			// log.Fatal(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.UnixConn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		utils.LogStderr("Could not read from connection")
		utils.LogStderr(err.Error())
		return
	}
	conn.Close()
}
