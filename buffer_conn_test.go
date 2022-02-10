package hybrid_udp_dtls_conn

import (
	"net"
	"reflect"
	"testing"

	"github.com/pion/udp"
)

func bufferClient(message []byte) {
	client, _ := net.Dial("udp", "127.0.0.1:8080")
	client.Write(message)
}

func TestRead(test *testing.T) {
	message := []byte{0x01, 0x01}
	listener, error := udp.Listen("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080})

	if error != nil {
		test.Error("ListenUDP error ", error)
	}

	defer listener.Close()
	go bufferClient(message)

	for {
		connection, error := listener.Accept()

		if error != nil {
			test.Error("Accept error ", error)
		}

		connectionBuffer := CreateBufferConn(connection)
		defer connectionBuffer.Close()
		buffer := make([]byte, 1024)
		length, error := connectionBuffer.Read(buffer)

		if error != nil {
			test.Error("Read error ", error)
		}

		test.Log("message", buffer[:length])

		if length != len(message) || !reflect.DeepEqual(message, buffer[:length]) {
			test.Error("messages not match ", message, buffer[:length])
		}

		break
	}
}

func TestReadBuffer(test *testing.T) {
	message := []byte{0x01, 0x01}
	toBuffer := []byte{0x02, 0x02, 0x02, 0x02}
	listener, error := udp.Listen("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080})

	if error != nil {
		test.Error("ListenUDP error ", error)
	}

	defer listener.Close()
	go bufferClient(message)

	for {
		connection, error := listener.Accept()

		if error != nil {
			test.Error("Accept error ", error)
		}

		connectionBuffer := CreateBufferConn(connection)
		connectionBuffer.SetBuffer(toBuffer)
		defer connectionBuffer.Close()
		buffer := make([]byte, 1024)
		length, error := connectionBuffer.Read(buffer)

		if error != nil {
			test.Error("Read error ", error)
		}

		test.Log("message", buffer[:length])

		if length != len(toBuffer) || !reflect.DeepEqual(toBuffer, buffer[:length]) {
			test.Error("messages not match ", toBuffer, buffer[:length])
		}

		break
	}
}
