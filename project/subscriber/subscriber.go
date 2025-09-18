package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/quic-go/quic-go"
	"jarkom.cs.ui.ac.id/h01/project/utils"
)

const (
	serverIP      = "" // IP catch-all
	serverPort    = "7101"
	serverType    = "udp4"
	appLayerProto = "lrt-jabodebek-2306207101"
	bufferSize    = 2048
)

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

	tlsConfig := &tls.Config{
		Certificates: utils.GenerateTLSSelfSignedCertificates(),
		NextProtos:   []string{appLayerProto},
	}

	listener, err := quic.Listen(socket, tlsConfig, &quic.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	//fmt.Printf("[quic] Listening QUIC connections on %s\n", listener.Addr())

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Println("connection accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn quic.Connection) {
	//fmt.Printf("[quic] Connected from %s\n", conn.RemoteAddr())

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			//log.Println("stream accept error:", err)
			return
		}
		go handleStream(stream)
	}
}

func handleStream(stream quic.Stream) {
	defer stream.Close()

	buffer := make([]byte, bufferSize)
	for {
		n, err := stream.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Println("stream read error:", err)
			}
			return
		}

		packet := utils.Decoder(buffer[:n])
		response := Handler(packet)

		fmt.Println(response) // print to PIDS display

		// Send ACK back
		ack := packet
		ack.IsAck = true
		ackBytes := utils.Encoder(ack)
		_, err = stream.Write(ackBytes)
		if err != nil {
			log.Println("stream write error:", err)
			return
		}
	}
}

// Generate message based on train arriving/departing
func Handler(packet utils.LRTPIDSPacket) string {
	if packet.IsTrainArriving {
		return fmt.Sprintf("Mohon perhatian, kereta tujuan %s akan tiba di Peron 1.", packet.Destination)
	}
	if packet.IsTrainDeparting {
		return fmt.Sprintf("Mohon perhatian, kereta tujuan %s akan diberangkatkan dari Peron 1.", packet.Destination)
	}
	return "Tidak ada pengumuman."
}
