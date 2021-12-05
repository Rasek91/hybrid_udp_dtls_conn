# hybrid_udp_dtls_conn

A net.Conn implementation which can change dtls.Conn from udp.Conn when a DTLS
ClientHello received. It should be used with udp.Listen.Accept() as underlying
connection for Buffer_Conn and the Buffer_coon should be used for Conn as
underlying connection. For implementation example check the test files.

## Usage

#### type Buffer_Conn

```go
type Buffer_Conn struct {
}
```

Buffer_Conn implements the net.Conn interface. It has an internal buffer and
when there is data in it the Read returns the data from the buffer and delete
the buffer instead of reading the underlying connection.

#### func  Create_Buffer_Conn

```go
func Create_Buffer_Conn(connection net.Conn) *Buffer_Conn
```

Create_Buffer_Conn returns a new Buffer_Conn using connection as the underlying
connection.

#### func (*Buffer_Conn) Close

```go
func (connection *Buffer_Conn) Close() error
```

Close the underlying connection.

#### func (*Buffer_Conn) Get_Buffer

```go
func (connection *Buffer_Conn) Get_Buffer() []byte
```

Get_Buffer read the data from the buffer and delete the data from it.

#### func (*Buffer_Conn) LocalAddr

```go
func (connection *Buffer_Conn) LocalAddr() net.Addr
```

LocalAddr of the underlying connection.

#### func (*Buffer_Conn) Read

```go
func (connection *Buffer_Conn) Read(buffer []byte) (int, error)
```

Read returns the data from the buffer and delete the buffer instead of reading
the underlying connection if the internal buffer has data.

#### func (*Buffer_Conn) RemoteAddr

```go
func (connection *Buffer_Conn) RemoteAddr() net.Addr
```

RemoteAddr of the underlying connection.

#### func (*Buffer_Conn) SetDeadline

```go
func (connection *Buffer_Conn) SetDeadline(time time.Time) error
```

SetDeadline to the underlying connection.

#### func (*Buffer_Conn) SetReadDeadline

```go
func (connection *Buffer_Conn) SetReadDeadline(time time.Time) error
```

SetReadDeadline to the underlying connection.

#### func (*Buffer_Conn) SetWriteDeadline

```go
func (connection *Buffer_Conn) SetWriteDeadline(time time.Time) error
```

SetWriteDeadline to the underlying connection.

#### func (*Buffer_Conn) Set_Buffer

```go
func (connection *Buffer_Conn) Set_Buffer(buffer []byte)
```

Set_Buffer the internal buffer will be equal to buffer. If there was data in the
internal buffer it will be overwritten.

#### func (*Buffer_Conn) Write

```go
func (connection *Buffer_Conn) Write(buffer []byte) (int, error)
```

Write to the underlying connection.

#### type Conn

```go
type Conn struct {
}
```

Conn implements the net.Conn interface. Which can change the underlying
connection from UDP to DTLS if a DTLS ClientHello received.

#### func  Create_Conn

```go
func Create_Conn(connection net.Conn, tls_config *dtls.Config) *Conn
```

Create_Conn returns a new Conn using connection converted to Buffer_Conn as the
underlying connection. The configuration config must be non-nil and must include
at least one certificate or else set GetCertificate, if TLS will be added to the
connection.

#### func (*Conn) Close

```go
func (connection *Conn) Close() error
```

Close the underlying connection.

#### func (*Conn) Get_TLS

```go
func (connection *Conn) Get_TLS() bool
```

Get_TLS returns true if the underlying connection is using TLS and false if not.

#### func (*Conn) Get_TLS_Config

```go
func (connection *Conn) Get_TLS_Config() *dtls.Config
```

Get_TLS_Config returns the TLS server configuration.

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

#### func (*Conn) SetWriteDeadline

```go
func (connection *Conn) SetWriteDeadline(time time.Time) error
```

SetWriteDeadline to the underlying connection.

#### func (*Conn) Set_TLS_Config

```go
func (connection *Conn) Set_TLS_Config(tls_config *dtls.Config)
```

Set_TLS_Config change the TLS server configuration. New connection will be not
generated if you change it and TLS is already in use.

#### func (*Conn) Write

```go
func (connection *Conn) Write(buffer []byte) (int, error)
```

Write to the underlying connection.
