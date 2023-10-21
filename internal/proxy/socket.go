package proxy

import (
  "encoding/json"
	"io"
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
		go handleClientConnection(conn)
	}
}

func handleClientConnection(conn *net.UnixConn) {
	pr, rw := io.Pipe()

	go func() {
		for {
			buf := make([]byte, utils.BufferSize)
			_, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					rw.Close()
					return
				}
			}
			rw.Write(buf)
		}
	}()

	jsonDecoder := json.NewDecoder(pr)

	for (true) {
		var cmd map[string]interface{}
		err := jsonDecoder.Decode(&cmd)
		if err != nil {
			if err == io.EOF {
				return
			}
			utils.LogStderr("Error decoding JSON")
			log.Fatal(err)
		}

		command := cmd["command"].(string)
		oneshot := cmd["oneshot"].(bool)

		switch command {
		case "echo":
			handleEchoCommand(conn, &cmd)
		}
		
		if oneshot {
			return
		}
	}
}

func handleEchoCommand(conn *net.UnixConn, cmd *map[string]interface{}) {
	jsonBytes, _ := json.Marshal(cmd)
	conn.Write(jsonBytes)
}
