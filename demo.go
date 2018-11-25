package main

import (
	"net"
	"fmt"
	"encoding/json"
)

const (
	AdbPort = "5037"
)

type AdbResp struct {
	Status  string
	Length  string
	Content string
}

func buildAdbConnection() net.Conn {
	conn, err := net.Dial("tcp", ":"+AdbPort)
	if err != nil {
		panic(err)
	}

	return conn
}

func encodeContent(content string) []byte {
	contentByte := []byte(content)
	contentLength := len(contentByte)
	contentLengthHex := []byte(fmt.Sprintf("%04X", contentLength))
	readyMsg := append(contentLengthHex, contentByte...)
	return readyMsg
}

func sendMsg(c net.Conn, content string) {
	readyMsg := encodeContent(content)
	c.Write(readyMsg)
}

func readData(c net.Conn, dataLen int) string {
	buf := make([]byte, dataLen)
	readOutLen, err := c.Read(buf)
	if err != nil || readOutLen == 0 {
		return ""
	}
	return string(buf[0:readOutLen])
}

func getCurrentDevices(c net.Conn) *AdbResp {
	sendMsg(c, "host:devices")
	status := readData(c, 4)
	dataLen := readData(c, 4)
	content := readData(c, 1024)

	newResp := new(AdbResp)
	newResp.Status = status
	newResp.Length = dataLen
	newResp.Content = content

	return newResp
}

func main() {
	adbConnection := buildAdbConnection()
	currentDevices := getCurrentDevices(adbConnection)

	// print result
	resultStr, err := json.Marshal(currentDevices)
	if err != nil {
		panic(err)
	}
	println(string(resultStr))

	// close socket after usage
	adbConnection.Close()
}
