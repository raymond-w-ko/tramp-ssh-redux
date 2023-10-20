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

var conn *net.UnixConn

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


func SetupConnectionToProxy() {
	var err error

	var addr *net.UnixAddr
	addr, err = net.ResolveUnixAddr("unix", utils.SocketFile)
	if err != nil {
		utils.FatalMessageAsJson("Could open addr for UNIX socket " + utils.SocketFile)
	}

	conn, err = net.DialUnix("unix", nil, addr)
	if err != nil {
		utils.FatalMessageAsJson("Could not dial UNIX socket " + utils.SocketFile)
	}
	defer conn.Close()

	cmd := readNextJsonObject()
	jsonBytes, _ := json.Marshal(cmd)
	_, err = conn.Write(jsonBytes)

	var resultBuffer bytes.Buffer
	for {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			if err == io.EOF {
				continue
			}
			utils.FatalMessageAsJson("Error reading byte when collecting JSON result")
		}
		resultBuffer.Write(b[:n])
		if b[n-1] == '}' {
			break
		}
	}
	jsonResult := readNextJsonObjectFromStdin(resultBuffer)
	utils.WriteJsonObjectToStdout(jsonResult)
}
