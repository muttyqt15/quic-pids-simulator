package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/quic-go/quic-go"
	"jarkom.cs.ui.ac.id/h01/project/utils"
)

const (
	serverIP      = "127.0.0.1"
	serverPort    = "54321"
	serverType    = "udp4"
	appLayerProto = "lrt-jabodebek-2306207101"
)

// Handler: process received packet and return message to display
func Handler(packet utils.LRTPIDSPacket) string {
	if packet.IsTrainArriving {
		return fmt.Sprintf("Mohon perhatian, kereta tujuan %s akan tiba di Peron 1.", packet.Destination)
	}
	if packet.IsTrainDeparting {
		return fmt.Sprintf("Mohon perhatian, kereta tujuan %s akan diberangkatkan dari Peron 1.", packet.Destination)
	}
	return "Tidak ada pengumuman."
}

func main() {
	localUdpAddress, err := net.ResolveUDPAddr(serverType, net.JoinHostPort(serverIP, serverPort))
	if err != nil {
		log.Fatalln(err)
	}

	socket, err := net.ListenUDP(serverType, localUdpAddress)
	if err != nil {
		log.Fatalln(err)
	}
	defer socket.Close()

	fmt.Printf("QUIC Server Socket Program Example in Go\n")
	fmt.Printf("[%s] Listening on %s\n", serverType, socket.LocalAddr())

	tlsConfig := &tls.Config{
		//Certificates: utils.GenerateTLSSelfSignedCertificates(),
		NextProtos: []string{appLayerProto},
	}

	listener, err := quic.Listen(socket, tlsConfig, &quic.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	fmt.Printf("[quic] Listening QUIC connections on %s\n", listener.Addr())

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn quic.Connection) {
	fmt.Printf("[quic] New connection from %s\n", conn.RemoteAddr())

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}
		go handleStream(conn, stream)
	}
}

func handleStream(conn quic.Connection, stream quic.Stream) {
	defer stream.Close()

	raw := make([]byte, 0)
	buf := make([]byte, 2048)

	n, err := stream.Read(buf)
	if err != nil && err != io.EOF {
		log.Println("Error reading stream:", err)
		return
	}
	raw = append(raw, buf[:n]...)

	packet := utils.Decoder(raw)

	// Process packet
	msg := Handler(packet)
	fmt.Println(msg)

	// Send ACK back
	ackPacket := packet
	ackPacket.IsAck = true
	ackPacket.TransactionId ^= 0xABCD
	ackBytes := utils.Encoder(ackPacket)

	_, err = stream.Write(ackBytes)
	if err != nil {
		log.Println("Error sending ACK:", err)
		return
	}
}
