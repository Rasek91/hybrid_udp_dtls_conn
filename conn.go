//A net.Conn implementation which can change dtls.Conn from udp.Conn when a DTLS ClientHello received.
//It should be used with udp.Listen.Accept() as underlying connection for BufferConn and the BufferConn should be used for Conn as underlying connection.
//For implementation example check the test files.
package hybrid_udp_dtls_conn

import (
	"net"
	"sync"
	"time"

	"github.com/pion/dtls"
)

//Conn implements the net.Conn interface.
//Which can change the underlying connection from UDP to DTLS if a DTLS ClientHello received.
type Conn struct {
	connection net.Conn
	lock       sync.RWMutex
	tlsConfig  *dtls.Config
	tls        bool
}

//Read from the underlying connection.
//If DTLS ClientHello received upgrade the connection to DTLS and read again for the actual message.
func (connection *Conn) Read(buffer []byte) (int, error) {
	var connectionBuffer *BufferConn
	var connectionTls net.Conn
	var error error
	length, error := connection.connection.Read(buffer)

	if buffer[0] == byte(0x16) && buffer[13] == byte(0x1) && (length-13) == int(int32(buffer[11])<<8+int32(buffer[12])) && (length-25) == int(int32(buffer[14])<<16+int32(buffer[15])<<8+int32(buffer[16])) {
		connectionBuffer = connection.connection.(*BufferConn)
		connectionBuffer.SetBuffer(buffer[:length])
		connection.lock.Lock()
		connectionTls, error = dtls.Server(connection.connection, connection.tlsConfig)

		if error != nil {
			return length, error
		}

		connection.connection = connectionTls
		connection.tls = true
		connection.lock.Unlock()
		length, error = connection.connection.Read(buffer)
	}

	return length, error
}

//Write to the underlying connection.
func (connection *Conn) Write(buffer []byte) (int, error) {
	connection.lock.Lock()
	defer connection.lock.Unlock()
	length, error := connection.connection.Write(buffer)
	return length, error
}

//Close the underlying connection.
func (connection *Conn) Close() error {
	error := connection.connection.Close()
	return error
}

//LocalAddr of the underlying connection.
func (connection *Conn) LocalAddr() net.Addr {
	address := connection.connection.LocalAddr()
	return address
}

//RemoteAddr of the underlying connection.
func (connection *Conn) RemoteAddr() net.Addr {
	address := connection.connection.RemoteAddr()
	return address
}

//SetDeadline to the underlying connection.
func (connection *Conn) SetDeadline(time time.Time) error {
	error := connection.connection.SetDeadline(time)
	return error
}

//SetReadDeadline to the underlying connection.
func (connection *Conn) SetReadDeadline(time time.Time) error {
	error := connection.connection.SetReadDeadline(time)
	return error
}

//SetWriteDeadline to the underlying connection.
func (connection *Conn) SetWriteDeadline(time time.Time) error {
	error := connection.connection.SetWriteDeadline(time)
	return error
}

//GetTls returns true if the underlying connection is using TLS and false if not.
func (connection *Conn) GetTls() bool {
	return connection.tls
}

//SetTlsConfig change the TLS server configuration.
//New connection will be not generated if you change it and TLS is already in use.
func (connection *Conn) SetTlsConfig(tlsConfig *dtls.Config) {
	connection.tlsConfig = tlsConfig
}

//GetTlsConfig returns the TLS server configuration.
func (connection *Conn) GetTlsConfig() *dtls.Config {
	return connection.tlsConfig
}

//New returns a new Conn using connection converted to BufferConn as the underlying connection.
//The configuration config must be non-nil and must include at least one certificate or else set GetCertificate, if TLS will be added to the connection.
func New(connection net.Conn, tlsConfig *dtls.Config) *Conn {
	return &Conn{connection: CreateBufferConn(connection), lock: sync.RWMutex{}, tlsConfig: tlsConfig}
}
