package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

const SocketFile = "./tramp-ssh-redux.sock"
const BufferSize = 8192

////////////////////////////////////////////////////////////////////////////////////////////////////

func LogStderr(v ...any) {
	for _, x := range v {
		switch y := x.(type) {
		case string:
			fmt.Fprintf(os.Stderr, "%s\n", y)
		case error:
			fmt.Fprintf(os.Stderr, "%s\n", y.Error())
		default:
			fmt.Fprintf(os.Stderr, "Unknown type\n")
		}
	}
}

func FatalStderr(v ...any) {
	LogStderr(v...)
	os.Exit(1)
}

func WriteJsonObjectToStdout(data map[string]interface{}) {
	jsonBytes, _ := json.MarshalIndent(data, "", "\t")
	os.Stdout.Write(jsonBytes)
}

func FatalMessageAsJson(v ...any) {
	var message []string
	for _, x := range v {
		switch y := x.(type) {
		case string:
			message = append(message, y)
		case error:
			message = append(message, y.Error())
		default:
		}
	}
	data := map[string]interface{}{
		"success": false,
		"error":   true,
		"message": message,
	}
	WriteJsonObjectToStdout(data)
	os.Exit(1)
}

func ChdirToSelfExecutablePath() {
	ex, err := os.Executable()
	if err != nil {
		LogStderr("Could not get self executable path")
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	if err := os.Chdir(exPath); err != nil {
		LogStderr("Could not change working directory to self executable path")
		log.Fatal(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

var pendingStdinBytes []byte = nil

func ReadStdinN(n uint64) []byte {
	bytes := make([]byte, 0, n)
	pendingLength := uint64(len(pendingStdinBytes))
	if pendingLength > 0 {
		if pendingLength <= n {
			bytes = append(bytes, pendingStdinBytes...)
			pendingStdinBytes = nil
			n -= pendingLength
		} else {
			a := pendingStdinBytes[:n]
			b := pendingStdinBytes[n:]
			bytes = append(bytes, a...)
			pendingStdinBytes = b
			return bytes
		}
	}
	if n == 0 {
		return bytes
	}

	for n > 0 {
		buf := make([]byte, BufferSize)
		numReadBytes, err := os.Stdin.Read(buf)
		m := uint64(numReadBytes)
		if err != nil {
			LogStderr("ReadStdinN: Could not read from stdin")
			LogStderr(err.Error())
			return bytes
		}
		buf = buf[:m]
		if m <= n {
			bytes = append(bytes, buf...)
			n -= m
		} else {
			a := buf[:n]
			b := buf[n:]
			bytes = append(bytes, a...)
			pendingStdinBytes = b
			return bytes
		}
	}

	return bytes
}

func ReadStdinUint32() uint32 {
	data := ReadStdinN(4)
	value := binary.BigEndian.Uint32(data)
	return value
}

func ReadStdinUint64() uint64 {
	data := ReadStdinN(8)
	value := binary.BigEndian.Uint64(data)
	fmt.Fprintf(os.Stderr, "ReadStdinUint64: %d\n", value)
	return value
}
