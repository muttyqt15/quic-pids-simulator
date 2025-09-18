package main

// CLIENT QUIC [kendali stasiun]
import (
	"fmt"

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

func main() {
	fmt.Println(serverIP)
	fmt.Println(serverPort)
	fmt.Println(serverType)
	fmt.Println(bufferSize)
	fmt.Println(appLayerProto)
	fmt.Println(sslKeyLogFileName)
	destination := "Harjamukti"
	packetTiba := utils.LRTPIDSPacket{
		LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
			TransactionId:     0x55,
			IsAck:             false,
			IsNewTrain:        false,
			IsUpdateTrain:     false,
			IsDeleteTrain:     false,
			IsTrainArriving:   true,
			IsTrainDeparting:  false,
			TrainNumber:       42,
			DestinationLength: uint8(len(destination)),
		},
		Destination: destination,
	}
	packetBerangkat := utils.LRTPIDSPacket{
		LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
			TransactionId:     0x55,
			IsAck:             false,
			IsNewTrain:        false,
			IsUpdateTrain:     false,
			IsDeleteTrain:     false,
			IsTrainArriving:   false,
			IsTrainDeparting:  true,
			TrainNumber:       42,
			DestinationLength: uint8(len(destination)),
		},
		Destination: destination,
	}
	result := utils.Encoder(packetTiba)
	result2 := utils.Encoder(packetBerangkat)
	fmt.Println(result)
	fmt.Println(result2)
	fmt.Println(utils.Decoder(result))

	fmt.Println("============= QUIC CONNECT LOGIC =============")

	//sslKeyLogFile, err := os.Create(sslKeyLogFileName)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer sslKeyLogFile.Close()
	//
	//fmt.Printf("QUIC Client Socket Program Example in Go\n")
	//
	//tlsConfig := &tls.Config{
	//	InsecureSkipVerify: true,
	//	NextProtos:         []string{appLayerProto},
	//	KeyLogWriter:       sslKeyLogFile,
	//}
	//connection, err := quic.DialAddr(context.Background(), net.JoinHostPort(serverIP, serverPort), tlsConfig, &quic.Config{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//defer connection.CloseWithError(0x0, "No Error")
	//
	//fmt.Printf("[quic] Dialling from %s to %s\n", connection.LocalAddr(), connection.RemoteAddr())
	//
	//fmt.Printf("[quic] Creating receive buffer of size %d\n", bufferSize)
	//receiveBuffer := make([]byte, bufferSize)
	//
	//fmt.Printf("[quic] Input message to be sent to server: ")
	//message, err := bufio.NewReader(os.Stdin).ReadString('\n')
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//var g sync.WaitGroup // ensures they run
	//for {                // unlimited streams
	//	g.Add(1)
	//	go func() {
	//		defer g.Done()
	//		// only 1 stream, but since in go routine multiple can run concurrently
	//		stream, err := connection.OpenStreamSync(context.Background())
	//		if err != nil {
	//			log.Fatalln(err)
	//		}
	//		fmt.Printf("[quic] Opened bidirectional stream %d to %s\n", stream.StreamID(), connection.RemoteAddr())
	//
	//		fmt.Printf("[quic] [Stream ID: %d] Sending message '%s'\n", stream.StreamID(), message)
	//		_, err = stream.Write([]byte(message))
	//		if err != nil {
	//			log.Fatalln(err)
	//		}
	//		fmt.Printf("[quic] [Stream ID: %d] Message sent\n", stream.StreamID())
	//
	//		receiveLength, err := stream.Read(receiveBuffer)
	//		if err != nil {
	//			log.Fatalln(err)
	//		}
	//		fmt.Printf("[quic] [Stream ID: %d] Received %d bytes of message from server\n", stream.StreamID(), receiveLength)
	//
	//		response := receiveBuffer[:receiveLength]
	//		fmt.Printf("[quic] [Stream ID: %d] Received message: '%s'\n", stream.StreamID(), response)
	//	}()
	//}
	//g.Wait()
}
