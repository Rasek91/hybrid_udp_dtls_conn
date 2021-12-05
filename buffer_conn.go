package hybrid_udp_dtls_conn

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

//Buffer_Conn implements the net.Conn interface.
//It has an internal buffer and when there is data in it the Read returns the data from the buffer and delete the buffer instead of reading the underlying connection.
type Buffer_Conn struct {
	connection net.Conn
	buffer     []byte
}

//Read returns the data from the buffer and delete the buffer instead of reading the underlying connection if the internal buffer has data.
func (connection *Buffer_Conn) Read(buffer []byte) (int, error) {
	if connection.buffer == nil {
		length, error := connection.connection.Read(buffer)
		log.Trace("get message ", buffer[:length])
		return length, error
	} else {
		log.Trace("message in buffer ", connection.buffer)
		length := copy(buffer, connection.read_buffer())
		return length, nil
	}
}

//Write to the underlying connection.
func (connection *Buffer_Conn) Write(buffer []byte) (int, error) {
	length, error := connection.connection.Write(buffer)
	log.Trace("sent message ", buffer[:length])
	return length, error
}

//Close the underlying connection.
func (connection *Buffer_Conn) Close() error {
	error := connection.connection.Close()
	log.Trace("connection close")
	return error
}

//LocalAddr of the underlying connection.
func (connection *Buffer_Conn) LocalAddr() net.Addr {
	address := connection.connection.LocalAddr()
	log.Trace("LocalAddr")
	return address
}

//RemoteAddr of the underlying connection.
func (connection *Buffer_Conn) RemoteAddr() net.Addr {
	address := connection.connection.RemoteAddr()
	log.Trace("RemoteAddr")
	return address
}

//SetDeadline to the underlying connection.
func (connection *Buffer_Conn) SetDeadline(time time.Time) error {
	error := connection.connection.SetDeadline(time)
	log.Trace("SetDeadline")
	return error
}

//SetReadDeadline to the underlying connection.
func (connection *Buffer_Conn) SetReadDeadline(time time.Time) error {
	error := connection.connection.SetReadDeadline(time)
	log.Trace("SetReadDeadline")
	return error
}

//SetWriteDeadline to the underlying connection.
func (connection *Buffer_Conn) SetWriteDeadline(time time.Time) error {
	error := connection.connection.SetWriteDeadline(time)
	log.Trace("SetWriteDeadline")
	return error
}

//Set_Buffer the internal buffer will be equal to buffer.
//If there was data in the internal buffer it will be overwritten.
func (connection *Buffer_Conn) Set_Buffer(buffer []byte) {
	connection.buffer = buffer
}

//Get_Buffer read the data from the buffer and delete the data from it.
func (connection *Buffer_Conn) Get_Buffer() []byte {
	return connection.read_buffer()
}

func (connection *Buffer_Conn) delete_buffer() {
	connection.buffer = nil
}

func (connection *Buffer_Conn) read_buffer() []byte {
	defer connection.delete_buffer()
	return connection.buffer
}

//Create_Buffer_Conn returns a new Buffer_Conn using connection as the underlying connection.
func Create_Buffer_Conn(connection net.Conn) *Buffer_Conn {
	return &Buffer_Conn{connection: connection}
}
