package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
)

// Tipe data komponen boleh diubah, namun variabelnya jangan diubah
type LRTPIDSPacketFixed struct {
	TransactionId     uint16
	IsAck             bool
	IsNewTrain        bool
	IsUpdateTrain     bool
	IsDeleteTrain     bool
	IsTrainArriving   bool
	IsTrainDeparting  bool
	TrainNumber       uint16
	DestinationLength uint8
}

type LRTPIDSPacket struct {
	LRTPIDSPacketFixed
	Destination string
}

func Encoder(packet LRTPIDSPacket) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, packet.LRTPIDSPacketFixed)
	if err != nil {
		log.Fatal(err)
	}
	// perlu ditangani secara terpisah, cuz gabisa langsung diconvert kalau string (non-fixed size)
	buffer.WriteString(packet.Destination)
	return buffer.Bytes()
}

func Decoder(rawMessage []byte) LRTPIDSPacket {
	var decodedPacketFixed LRTPIDSPacketFixed
	bytesReader := bytes.NewReader(rawMessage)
	err := binary.Read(bytesReader, binary.BigEndian, &decodedPacketFixed)
	if err != nil {
		log.Fatal(err)
	}
	destionationBytes := make([]byte, decodedPacketFixed.DestinationLength)
	n, err := bytesReader.Read(destionationBytes)
	if err != nil {
		log.Fatal(err)
	}
	if n != int(decodedPacketFixed.DestinationLength) {
		log.Fatalln(fmt.Errorf("Expected %d bytes but got %d", decodedPacketFixed.DestinationLength, n))
	}

	decodedPacket := LRTPIDSPacket{
		LRTPIDSPacketFixed: decodedPacketFixed,
		Destination:        string(destionationBytes),
	}

	return decodedPacket
}

func GenerateTLSSelfSignedCertificates() []tls.Certificate {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		log.Fatalln(err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatalln(err)
	}
	return []tls.Certificate{tlsCert}
}
