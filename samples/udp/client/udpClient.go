package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	serverIP   = "52.73.42.39"
	serverPort = "7101"
	serverType = "udp4"
	bufferSize = 2048
)

func main() {
	remoteUdpAddress, err := net.ResolveUDPAddr(serverType, net.JoinHostPort(serverIP, serverPort))
	if err != nil {
		log.Fatalln(err)
	}
	socket, err := net.DialUDP(serverType, nil, remoteUdpAddress)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("UDP Client Socket Program Example in Go\n")
	fmt.Printf("[%s] Dialling from %s to %s\n", serverType, socket.LocalAddr(), socket.RemoteAddr())

	defer socket.Close()

	fmt.Printf("[%s] Creating receive buffer of size %d\n", serverType, bufferSize)
	receiveBuffer := make([]byte, bufferSize)

	fmt.Printf("[%s] Input message to be sent to server: ", serverType)
	message, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("[%s] Sending message '%s' to server\n", serverType, message)
	_, err = socket.Write([]byte(message))
	if err != nil {
		log.Fatalln(err)
	}

	receiveLength, err := socket.Read(receiveBuffer) // blocking func until message is recieved
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("[%s] Received %d bytes of message from server\n", serverType, receiveLength)

	response := string(receiveBuffer[:receiveLength])
	fmt.Printf("[%s] Response from server: %s\n", serverType, response)

}
