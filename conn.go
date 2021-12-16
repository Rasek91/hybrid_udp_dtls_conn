//A net.Conn implementation which can change dtls.Conn from udp.Conn when a DTLS ClientHello received.
//It should be used with udp.Listen.Accept() as underlying connection for Buffer_Conn and the Buffer_coon should be used for Conn as underlying connection.
//For implementation example check the test files.
package hybrid_udp_dtls_conn

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/pion/dtls"

	log "github.com/sirupsen/logrus"
)

//Conn implements the net.Conn interface.
//Which can change the underlying connection from UDP to DTLS if a DTLS ClientHello received.
type Conn struct {
	connection    net.Conn
	lock          sync.RWMutex
	tls_config    *dtls.Config
	tls           bool
	read_sync     *bool
	write_sync    *bool
	read_channel  chan result
	write_channel chan result
	read_context  context.Context
	write_context context.Context
	read_cancel   context.CancelFunc
	write_cancel  context.CancelFunc
}

type result struct {
	error error
	length int
}

type my_error struct {
	Op string
	Net string
	Source net.Addr
	Addr net.Addr
	Err error
}

func (error *my_error) Error() string {
	if error == nil {
		return "<nil>"
	}
	string := error.Op
	if error.Net != "" {
		string += " " + error.Net
	}
	if error.Source != nil {
		string += " " + error.Source.String()
	}
	if error.Addr != nil {
		if error.Source != nil {
			string += "->"
		} else {
			string += " "
		}
		string += error.Addr.String()
	}
	string += ": " + error.Err.Error()
	return string
}

func (error *my_error) Timeout() bool {
	return true
}

//Read from the underlying connection.
//If DTLS ClientHello received upgrade the connection to DTLS and read again for the actual message.
func (connection *Conn) Read(buffer []byte) (int, error) {
	if !*connection.read_sync {
		go connection.read(buffer)
	}

	if connection.read_cancel != nil {
		defer connection.read_cancel()
	}

	select {
	case <- connection.read_context.Done():
		if len(connection.read_channel) != 0 {
			result := <- connection.read_channel
			return result.length, nil
		} else {
			return 0, &my_error{Op: "read", Net: "udp", Source: connection.connection.LocalAddr(), Addr: connection.connection.RemoteAddr(), Err: errors.New("timeout")}
		}
	case result := <- connection.read_channel:
		return result.length, result.error
	}
}

func (connection *Conn) read(buffer []byte) {
	*connection.read_sync = true
	var connection_buffer *Buffer_Conn
	var connection_tls net.Conn
	var error error
	length, error := connection.connection.Read(buffer)

	if buffer[0] == byte(0x16) && buffer[13] == byte(0x1) && (length-13) == int(int32(buffer[11])<<8+int32(buffer[12])) && (length-25) == int(int32(buffer[14])<<16+int32(buffer[15])<<8+int32(buffer[16])) {
		log.Trace("Client Hello received")
		connection_buffer = connection.connection.(*Buffer_Conn)
		connection_buffer.Set_Buffer(buffer[:length])
		connection.lock.Lock()
		connection_tls, error = dtls.Server(connection.connection, connection.tls_config)

		if error != nil {
			*connection.read_sync = false
			connection.read_channel <- result{length: length, error: error}
		}

		connection.connection = connection_tls
		log.Trace("TLS Handshake was successful")
		connection.tls = true
		connection.lock.Unlock()
		length, error = connection.connection.Read(buffer)
	}

	log.Trace("get message ", buffer[:length])
	*connection.read_sync = false
	result := result{length: length, error: error}
	connection.read_channel <- result
}

//Write to the underlying connection.
func (connection *Conn) Write(buffer []byte) (int, error) {
	if !*connection.write_sync {
		go connection.write(buffer)
	}

	if connection.write_cancel != nil {
		defer connection.write_cancel()
	}

	select {
	case <- connection.write_context.Done():
		if len(connection.write_channel) != 0 {
			result := <- connection.write_channel
			return result.length, result.error
		} else {
			return 0, &my_error{Op: "write", Net: "udp", Source: connection.connection.LocalAddr(), Addr: connection.connection.RemoteAddr(), Err: errors.New("timeout")}
		}
	case result := <- connection.write_channel:
		return result.length, result.error
	}
}

func (connection *Conn) write(buffer []byte) {
	*connection.write_sync = true
	connection.lock.Lock()
	defer connection.lock.Unlock()
	length, error := connection.connection.Write(buffer)
	log.Trace("sent message ", buffer[:length])
	result := result{length: length, error: error}
	connection.write_channel <- result
	*connection.write_sync = false
}

//Close the underlying connection.
func (connection *Conn) Close() error {
	error := connection.connection.Close()
	close(connection.read_channel)
	close(connection.write_channel)
	log.Trace("connection close")
	return error
}

//LocalAddr of the underlying connection.
func (connection *Conn) LocalAddr() net.Addr {
	address := connection.connection.LocalAddr()
	log.Trace("LocalAddr")
	return address
}

//RemoteAddr of the underlying connection.
func (connection *Conn) RemoteAddr() net.Addr {
	address := connection.connection.RemoteAddr()
	log.Trace("RemoteAddr")
	return address
}

//SetDeadline to the underlying connection.
func (connection *Conn) SetDeadline(time time.Time) error {
	connection.read_context, connection.read_cancel = context.WithDeadline(context.Background(), time)
	connection.write_context, connection.write_cancel = context.WithDeadline(context.Background(), time)
	log.Trace("SetDeadline")
	return nil
}

//SetReadDeadline to the underlying connection.
func (connection *Conn) SetReadDeadline(time time.Time) error {
	connection.read_context, connection.read_cancel = context.WithDeadline(context.Background(), time)
	log.Trace("SetReadDeadline")
	return nil
}

//SetWriteDeadline to the underlying connection.
func (connection *Conn) SetWriteDeadline(time time.Time) error {
	connection.write_context, connection.write_cancel = context.WithDeadline(context.Background(), time)
	log.Trace("SetWriteDeadline")
	return nil
}

//Get_TLS returns true if the underlying connection is using TLS and false if not.
func (connection *Conn) Get_TLS() bool {
	return connection.tls
}

//Set_TLS_Config change the TLS server configuration.
//New connection will be not generated if you change it and TLS is already in use.
func (connection *Conn) Set_TLS_Config(tls_config *dtls.Config) {
	connection.tls_config = tls_config
}

//Get_TLS_Config returns the TLS server configuration.
func (connection *Conn) Get_TLS_Config() *dtls.Config {
	return connection.tls_config
}

//Create_Conn returns a new Conn using connection converted to Buffer_Conn as the underlying connection.
//The configuration config must be non-nil and must include at least one certificate or else set GetCertificate, if TLS will be added to the connection.
func Create_Conn(connection net.Conn, tls_config *dtls.Config) *Conn {
	return &Conn{connection: Create_Buffer_Conn(connection), lock: sync.RWMutex{}, tls_config: tls_config, read_channel: make(chan result, 1), write_channel: make(chan result, 1), read_context: context.Background(), write_context: context.Background(), read_sync: new(bool), write_sync: new(bool)}
}
