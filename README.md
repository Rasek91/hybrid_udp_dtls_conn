# hybrid_udp_dtls_conn

A net.Conn implementation which can change dtls.Conn from udp.Conn when a DTLS
ClientHello received. It should be used with udp.Listen.Accept() as underlying
connection for BufferConn and the BufferConn should be used for Conn as
underlying connection. For implementation example check the test files.

## Usage

#### type BufferConn

```go
type BufferConn struct {
}
```

BufferConn implements the net.Conn interface. It has an internal buffer and when
there is data in it the Read returns the data from the buffer and delete the
buffer instead of reading the underlying connection.

#### func  CreateBufferConn

```go
func CreateBufferConn(connection net.Conn) *BufferConn
```

CreateBufferConn returns a new BufferConn using connection as the underlying
connection.

#### func (*BufferConn) Close

```go
func (connection *BufferConn) Close() error
```

Close the underlying connection.

#### func (*BufferConn) GetBuffer

```go
func (connection *BufferConn) GetBuffer() []byte
```

GetBuffer read the data from the buffer and delete the data from it.

#### func (*BufferConn) LocalAddr

```go
func (connection *BufferConn) LocalAddr() net.Addr
```

LocalAddr of the underlying connection.

#### func (*BufferConn) Read

```go
func (connection *BufferConn) Read(buffer []byte) (int, error)
```

Read returns the data from the buffer and delete the buffer instead of reading
the underlying connection if the internal buffer has data.

#### func (*BufferConn) RemoteAddr

```go
func (connection *BufferConn) RemoteAddr() net.Addr
```

RemoteAddr of the underlying connection.

#### func (*BufferConn) SetBuffer

```go
func (connection *BufferConn) SetBuffer(buffer []byte)
```

SetBuffer the internal buffer will be equal to buffer. If there was data in the
internal buffer it will be overwritten.

#### func (*BufferConn) SetDeadline

```go
func (connection *BufferConn) SetDeadline(time time.Time) error
```

SetDeadline to the underlying connection.

#### func (*BufferConn) SetReadDeadline

```go
func (connection *BufferConn) SetReadDeadline(time time.Time) error
```

SetReadDeadline to the underlying connection.

#### func (*BufferConn) SetWriteDeadline

```go
func (connection *BufferConn) SetWriteDeadline(time time.Time) error
```

SetWriteDeadline to the underlying connection.

#### func (*BufferConn) Write

```go
func (connection *BufferConn) Write(buffer []byte) (int, error)
```

Write to the underlying connection.

#### type Conn

```go
type Conn struct {
}
```

Conn implements the net.Conn interface. Which can change the underlying
connection from UDP to DTLS if a DTLS ClientHello received.

#### func  New

```go
func New(connection net.Conn, tlsConfig *dtls.Config) *Conn
```

New returns a new Conn using connection converted to BufferConn as the
underlying connection. The configuration config must be non-nil and must include
at least one certificate or else set GetCertificate, if TLS will be added to the
connection.

#### func (*Conn) Close

```go
func (connection *Conn) Close() error
```

Close the underlying connection.

#### func (*Conn) GetTls

```go
func (connection *Conn) GetTls() bool
```

GetTls returns true if the underlying connection is using TLS and false if not.

#### func (*Conn) GetTlsConfig

```go
func (connection *Conn) GetTlsConfig() *dtls.Config
```

GetTlsConfig returns the TLS server configuration.

#### func (*Conn) LocalAddr

```go
func (connection *Conn) LocalAddr() net.Addr
```

LocalAddr of the underlying connection.

#### func (*Conn) Read

```go
func (connection *Conn) Read(buffer []byte) (int, error)
```

Read from the underlying connection. If DTLS ClientHello received upgrade the
connection to DTLS and read again for the actual message.

#### func (*Conn) RemoteAddr

```go
func (connection *Conn) RemoteAddr() net.Addr
```

RemoteAddr of the underlying connection.

#### func (*Conn) SetDeadline

```go
func (connection *Conn) SetDeadline(time time.Time) error
```

SetDeadline to the underlying connection.

#### func (*Conn) SetReadDeadline

```go
func (connection *Conn) SetReadDeadline(time time.Time) error
```

SetReadDeadline to the underlying connection.

#### func (*Conn) SetTlsConfig

```go
func (connection *Conn) SetTlsConfig(tlsConfig *dtls.Config)
```

SetTlsConfig change the TLS server configuration. New connection will be not
generated if you change it and TLS is already in use.

#### func (*Conn) SetWriteDeadline

```go
func (connection *Conn) SetWriteDeadline(time time.Time) error
```

SetWriteDeadline to the underlying connection.

#### func (*Conn) Write

```go
func (connection *Conn) Write(buffer []byte) (int, error)
```

Write to the underlying connection.
