package main

// CLIENT QUIC [kendali stasiun]
import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/quic-go/quic-go"
	"jarkom.cs.ui.ac.id/h01/project/utils"
)

const (
	serverIP      = "127.0.0.1"
	serverPort    = "54321"
	bufferSize    = 2048
	appLayerProto = "lrt-jabodebek-2306207101"
)

func sendPacket(connection quic.Connection, packet utils.LRTPIDSPacket) {
	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer stream.Close()

	data := utils.Encoder(packet)

	// Send packet
	_, err = stream.Write(data)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[client] Sent packet (TransactionId: 0x%X, TrainNumber: %d)\n", packet.TransactionId, packet.TrainNumber)

	// Receive ACK
	receiveBuffer := make([]byte, bufferSize)
	n, err := stream.Read(receiveBuffer)
	if err != nil {
		log.Fatalln(err)
	}
	ackPacket := utils.Decoder(receiveBuffer[:n])
	fmt.Printf("[client] Received ACK (TransactionId: 0x%X, IsAck: %v)\n", ackPacket.TransactionId, ackPacket.IsAck)
}

func main() {
	sslKeyLogFile, err := os.Create("C:\\Users\\hayay\\Downloads\\Misc\\ssl-key.log")
	if err != nil {
		log.Fatalln(err)
	}
	defer sslKeyLogFile.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{appLayerProto},
		KeyLogWriter:       sslKeyLogFile,
	}

	connection, err := quic.DialAddr(context.Background(), net.JoinHostPort(serverIP, serverPort), tlsConfig, &quic.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	defer connection.CloseWithError(0x0, "No Error")

	destination := "Harjamukti"
	packetTiba := utils.LRTPIDSPacket{
		LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
			TransactionId:     0x55,
			IsTrainArriving:   true,
			TrainNumber:       42,
			DestinationLength: uint8(len(destination)),
		},
		Destination: destination,
	}

	packetBerangkat := utils.LRTPIDSPacket{
		LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
			TransactionId:     0x55,
			IsTrainDeparting:  true,
			TrainNumber:       42,
			DestinationLength: uint8(len(destination)),
		},
		Destination: destination,
	}

	// Send arrival packet
	sendPacket(connection, packetTiba)

	// Send departure packet
	sendPacket(connection, packetBerangkat)
}
