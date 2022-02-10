package hybrid_udp_dtls_conn

import (
	"net"
	"time"
)

//BufferConn implements the net.Conn interface.
//It has an internal buffer and when there is data in it the Read returns the data from the buffer and delete the buffer instead of reading the underlying connection.
type BufferConn struct {
	connection net.Conn
	buffer     []byte
}

//Read returns the data from the buffer and delete the buffer instead of reading the underlying connection if the internal buffer has data.
func (connection *BufferConn) Read(buffer []byte) (int, error) {
	if connection.buffer == nil {
		length, error := connection.connection.Read(buffer)
		return length, error
	} else {
		length := copy(buffer, connection.readBuffer())
		return length, nil
	}
}

//Write to the underlying connection.
func (connection *BufferConn) Write(buffer []byte) (int, error) {
	length, error := connection.connection.Write(buffer)
	return length, error
}

//Close the underlying connection.
func (connection *BufferConn) Close() error {
	error := connection.connection.Close()
	return error
}

//LocalAddr of the underlying connection.
func (connection *BufferConn) LocalAddr() net.Addr {
	address := connection.connection.LocalAddr()
	return address
}

//RemoteAddr of the underlying connection.
func (connection *BufferConn) RemoteAddr() net.Addr {
	address := connection.connection.RemoteAddr()
	return address
}

//SetDeadline to the underlying connection.
func (connection *BufferConn) SetDeadline(time time.Time) error {
	error := connection.connection.SetDeadline(time)
	return error
}

//SetReadDeadline to the underlying connection.
func (connection *BufferConn) SetReadDeadline(time time.Time) error {
	error := connection.connection.SetReadDeadline(time)
	return error
}

//SetWriteDeadline to the underlying connection.
func (connection *BufferConn) SetWriteDeadline(time time.Time) error {
	error := connection.connection.SetWriteDeadline(time)
	return error
}

//SetBuffer the internal buffer will be equal to buffer.
//If there was data in the internal buffer it will be overwritten.
func (connection *BufferConn) SetBuffer(buffer []byte) {
	connection.buffer = buffer
}

//GetBuffer read the data from the buffer and delete the data from it.
func (connection *BufferConn) GetBuffer() []byte {
	return connection.readBuffer()
}

func (connection *BufferConn) deleteBuffer() {
	connection.buffer = nil
}

func (connection *BufferConn) readBuffer() []byte {
	defer connection.deleteBuffer()
	return connection.buffer
}

//CreateBufferConn returns a new BufferConn using connection as the underlying connection.
func CreateBufferConn(connection net.Conn) *BufferConn {
	return &BufferConn{connection: connection}
}
