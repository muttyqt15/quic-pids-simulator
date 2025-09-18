package main

// Kendali Stasiun (Client, VM1)
import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"jarkom.cs.ui.ac.id/h01/project/utils"
)

const (
	serverIP      = "52.73.42.39" // IP Private VM1
	serverPort    = "7101"
	appLayerProto = "lrt-jabodebek-2306207101"
)

func main() {
	sslKeyLogFile, err := os.Create("C:\\Users\\hayay\\Downloads\\Misc\\ssl-key.log")
	if err != nil {
		log.Fatal(err)
	}
	defer sslKeyLogFile.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{appLayerProto},
		KeyLogWriter:       sslKeyLogFile,
	}

	conn, err := quic.DialAddr(context.Background(), net.JoinHostPort(serverIP, serverPort), tlsConfig, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.CloseWithError(0, "bye")

	destination := "Harjamukti"

	packets := []utils.LRTPIDSPacket{
		{
			LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
				TransactionId:     0x55,
				IsTrainArriving:   true,
				TrainNumber:       42,
				DestinationLength: uint8(len(destination)),
			},
			Destination: destination,
		},
		{
			LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
				TransactionId:     0x55,
				IsTrainDeparting:  true,
				TrainNumber:       42,
				DestinationLength: uint8(len(destination)),
			},
			Destination: destination,
		},
	}

	for _, pkt := range packets {
		go sendPacket(conn, pkt)
	}

	// Wait a little to allow goroutines to finish
	time.Sleep(3 * time.Second)
}

func sendPacket(conn quic.Connection, pkt utils.LRTPIDSPacket) {
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	packetBytes := utils.Encoder(pkt)

	_, err = stream.Write(packetBytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[client] Sent packet (TransactionId: 0x%x, TrainNumber: %d)\n", pkt.TransactionId, pkt.TrainNumber)

	// wait for ACK
	buffer := make([]byte, 1024)
	n, err := stream.Read(buffer)
	if err != nil {
		log.Println("[client] read error:", err)
		return
	}

	ack := utils.Decoder(buffer[:n])
	if ack.IsAck {
		fmt.Printf("[client] Received ACK for TransactionId: 0x%x\n", ack.TransactionId)
	}
}
