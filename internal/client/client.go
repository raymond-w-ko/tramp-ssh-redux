package client

import (
	// "bufio"
	"bytes"
	"encoding/json"
	"io"
	// "log"
	"net"
	"os"

	"github.com/raymond-w-ko/tramp-ssh-redux/pkg/utils"
)

func readNextJsonObjectFromStdin(inBuffer bytes.Buffer) map[string]interface{} {
	var outBuffer bytes.Buffer
	for {
		b, err := inBuffer.ReadByte()
		if err != nil {
			if err == io.EOF {
				continue
			}
			utils.FatalMessageAsJson("Error reading byte", err)
		}

		outBuffer.WriteByte(b)

		if b == '}' {
			var data map[string]interface{}
			if err := json.Unmarshal(outBuffer.Bytes(), &data); err != nil {
				utils.FatalMessageAsJson("Error decoding JSON:", err)
			}
			outBuffer.Reset()
			return data
		}
	}
}

func readNextJsonObject() map[string]interface{} {
	var data map[string]interface{}
	bites := []byte(os.Args[1])
	err := json.Unmarshal(bites, &data)
	if err != nil {
		utils.FatalMessageAsJson("Error decoding JSON:", err)
	}
	return data
}

var proxyConn *net.UnixConn

func SetupConnectionToProxy() {
	var err error

	var addr *net.UnixAddr
	addr, err = net.ResolveUnixAddr("unix", utils.SocketFile)
	if err != nil {
		utils.FatalMessageAsJson("Could open addr for UNIX socket " + utils.SocketFile)
	}

	proxyConn, err = net.DialUnix("unix", nil, addr)
	if err != nil {
		utils.FatalMessageAsJson("Could not dial UNIX socket " + utils.SocketFile)
	}
}

func ReadInitialCommandFromArgs() map[string]interface{} {
	pr, rw := io.Pipe()

	go func() {
		jsonBytes := []byte(os.Args[1])
		rw.Write(jsonBytes)
		rw.Close()
	}()

	jsonDecoder := json.NewDecoder(pr)
	cmd := make(map[string]interface{})
	err := jsonDecoder.Decode(&cmd)
	if err != nil {
		utils.FatalMessageAsJson("Error decoding JSON:", err)
	}
	return cmd
}

func SendCommandToProxyAndWriteOutputToStdout(cmd map[string]interface{}) {
	jsonInputBytes, _ := json.Marshal(cmd)
	proxyConn.Write(jsonInputBytes)

	pr, rw := io.Pipe()
	oneshot := cmd["oneshot"].(bool)

	go func() {
		for {
			buf := make([]byte, utils.BufferSize)
			_, err := proxyConn.Read(buf)
			if err != nil {
				if err == io.EOF {
					rw.Close()
					return
				}
			}
			rw.Write(buf)
		}
	}()

	if (oneshot) {
		jsonDecoder := json.NewDecoder(pr)
		var jsonResult map[string]interface{}
		err := jsonDecoder.Decode(&jsonResult)
		if err != nil {
			utils.FatalMessageAsJson("Error decoding oneshot result JSON:", err)
		}
		utils.WriteJsonObjectToStdout(jsonResult)
	}
}
