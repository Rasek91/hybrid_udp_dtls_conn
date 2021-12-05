package hybrid_udp_dtls_conn

import (
	"net"
	"reflect"
	"testing"

	"github.com/pion/udp"
)

func buffer_client(message []byte) {
	client, _ := net.Dial("udp", "127.0.0.1:8080")
	client.Write(message)
}

func Test_Read(test *testing.T) {
	message := []byte{0x01, 0x01}
	listener, error := udp.Listen("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080})

	if error != nil {
		test.Error("ListenUDP error ", error)
	}

	defer listener.Close()
	go buffer_client(message)

	for {
		connection, error := listener.Accept()

		if error != nil {
			test.Error("Accept error ", error)
		}

		connection_buffer := Create_Buffer_Conn(connection)
		defer connection_buffer.Close()
		buffer := make([]byte, 1024)
		length, error := connection_buffer.Read(buffer)

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

func Test_Read_Buffer(test *testing.T) {
	message := []byte{0x01, 0x01}
	to_buffer := []byte{0x02, 0x02, 0x02, 0x02}
	listener, error := udp.Listen("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080})

	if error != nil {
		test.Error("ListenUDP error ", error)
	}

	defer listener.Close()
	go buffer_client(message)

	for {
		connection, error := listener.Accept()

		if error != nil {
			test.Error("Accept error ", error)
		}

		connection_buffer := Create_Buffer_Conn(connection)
		connection_buffer.Set_Buffer(to_buffer)
		defer connection_buffer.Close()
		buffer := make([]byte, 1024)
		length, error := connection_buffer.Read(buffer)

		if error != nil {
			test.Error("Read error ", error)
		}

		test.Log("message", buffer[:length])

		if length != len(to_buffer) || !reflect.DeepEqual(to_buffer, buffer[:length]) {
			test.Error("messages not match ", to_buffer, buffer[:length])
		}

		break
	}
}
