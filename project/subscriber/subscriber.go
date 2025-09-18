package main

// SERVER QUIC [layar PIDS]
import (
	"context"
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
	serverIP          = "127.0.0.1"
	serverPort        = "54321"
	serverType        = "udp4"
	bufferSize        = 2048
	appLayerProto     = "lrt-jabodebek-2306207101"
	sslKeyLogFileName = "C:\\Users\\hayay\\Downloads\\Misc\\ssl-key.log"
)

func Handler(packet utils.LRTPIDSPacket) string {
	// Fungsi untuk menangani objek pesan harus ditempatkan di fungsi Handler pada
	// subscriber.go. Fungsi ini hanya ditugaskan untuk menerima masukan berupa objek yang
	// sudah di-decode dan mengembalikan pesan yang harus ditampikan sesuai konten informasi
	// yang diterima.
	return ""
}

func main() {
	//localUdpAddress, err := net.ResolveUDPAddr(serverType, net.JoinHostPort(serverIP, serverPort))
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//socket, err := net.ListenUDP(serverType, localUdpAddress)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//defer func(socket *net.UDPConn) {
	//	err := socket.Close()
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//}(socket)
	//
	//fmt.Printf("QUIC Server Socket Program Example in Go\n")
	//fmt.Printf("[%s] Preparing UDP listening socket on %s\n", serverType, socket.LocalAddr())
	//
	//tlsConfig := &tls.Config{
	//	//Certificates: utils.GenerateTLSSelfSignedCertificates(), // self signed
	//	NextProtos: []string{appLayerProto},
	//}
	//listener, err := quic.Listen(socket, tlsConfig, &quic.Config{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//defer listener.Close()
	//
	//fmt.Printf("[quic] Listening QUIC connections on %s\n", listener.Addr())
	//
	//for {
	//	connection, err := listener.Accept(context.Background())
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	go handleConnection(connection)
	//}
	fmt.Println(serverIP)
	fmt.Println(serverPort)
	fmt.Println(serverType)
	fmt.Println(bufferSize)
	fmt.Println(appLayerProto)
	fmt.Println(sslKeyLogFileName)
	destination := "Dukuh Atas"
	packet := utils.LRTPIDSPacket{
		LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
			TransactionId:     0x55,
			IsAck:             false,
			IsNewTrain:        false,
			IsUpdateTrain:     false,
			IsDeleteTrain:     false,
			IsTrainArriving:   false,
			IsTrainDeparting:  true,
			TrainNumber:       1000,
			DestinationLength: uint8(len(destination)),
		},
		Destination: destination,
	}
	result := utils.Encoder(packet)
	fmt.Println(result)
	fmt.Println(utils.Decoder(result))
}

func handleConnection(connection quic.Connection) {
	fmt.Printf("[quic] Receiving connection from %s\n", connection.RemoteAddr())

	var g sync.WaitGroup
	for {
		g.Add(1)
		go func() {
			defer g.Done()
			stream, err := connection.AcceptStream(context.Background())
			if err != nil {
				log.Fatalln(err)
			}
			go handleStream(connection.RemoteAddr(), stream)
		}()
	}
	g.Wait()
}

func handleStream(clientAddress net.Addr, stream quic.Stream) {
	fmt.Printf("[quic] [Client: %s] Receive stream open request with ID %d\n", clientAddress, stream.StreamID())

	_, err := io.Copy(logicProcessorAndWriter{stream}, stream)
	if err != nil {
		fmt.Println(err)
	}
}

type logicProcessorAndWriter struct{ io.Writer }

func (lp logicProcessorAndWriter) Write(receivedMessageRaw []byte) (int, error) {

	receivedMessage := string(receivedMessageRaw)
	fmt.Printf("[quic] Receive message: %s\n", receivedMessage)

	response := strings.ToUpper(receivedMessage)
	writeLength, err := lp.Writer.Write([]byte(response))

	fmt.Printf("[quic] Send message: %s\n", response)

	return writeLength, err
}
