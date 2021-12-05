package hybrid_udp_dtls_conn

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/pion/dtls"
	"github.com/pion/udp"
)

func udp_client(message []byte) {
	client, _ := net.Dial("udp", "127.0.0.1:8080")
	client.Write(message)
}

func tls_client(message []byte, test *testing.T) {
	client_udp, _ := net.Dial("udp", "127.0.0.1:8080")
	client, error := dtls.Client(client_udp, &dtls.Config{InsecureSkipVerify: true})

	if error != nil {
		test.Error("DTLS client error ", error)
	}

	client.Write(message)
}

func Test_UDP(test *testing.T) {
	message := []byte{0x01, 0x01}
	listener, error := udp.Listen("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080})

	if error != nil {
		test.Error("ListenUDP error ", error)
	}

	defer listener.Close()
	go udp_client(message)

	for {
		connection, error := listener.Accept()

		if error != nil {
			test.Error("Accept error ", error)
		}

		connection_hybrid := Create_Conn(connection, nil)
		defer connection_hybrid.Close()
		buffer := make([]byte, 1024)
		length, error := connection_hybrid.Read(buffer)

		if error != nil {
			test.Error("Read error ", error)
		}

		test.Log("message", buffer[:length])

		if connection_hybrid.Get_TLS() {
			test.Error("Get_TLS error")
		}

		if length != len(message) || !reflect.DeepEqual(message, buffer[:length]) {
			test.Error("messages not match ", message, buffer[:length])
		}

		break
	}
}

func Test_TLS(test *testing.T) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	ca_private_key, error := rsa.GenerateKey(rand.Reader, 4096)
	if error != nil {
		test.Error("CA Private Key error ", error)
	}

	ca_bytes, error := x509.CreateCertificate(rand.Reader, ca, ca, &ca_private_key.PublicKey, ca_private_key)
	if error != nil {
		test.Error("CA Bytes error ", error)
	}

	ca_pem := new(bytes.Buffer)
	pem.Encode(ca_pem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca_bytes,
	})

	ca_private_key_pem := new(bytes.Buffer)
	pem.Encode(ca_private_key_pem, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(ca_private_key),
	})

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	cert_private_key, error := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if error != nil {
		test.Error("Cert Private Key error ", error)
	}

	cert_bytes, error := x509.CreateCertificate(rand.Reader, cert, ca, &cert_private_key.PublicKey, ca_private_key)
	if error != nil {
		test.Error("Cert Bytes error ", error)
	}

	certificate, error := x509.ParseCertificate(cert_bytes)
	if error != nil {
		test.Error("Parse Cert Bytes error ", error)
	}


	server_tls_config := &dtls.Config{
		Certificate: certificate,
		PrivateKey: cert_private_key,
	}
	test.Log("TLS Cert generated")
	message := []byte{0x01, 0x01}
	listener, error := udp.Listen("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080})

	if error != nil {
		test.Error("ListenUDP error ", error)
	}

	defer listener.Close()
	go tls_client(message, test)

	for {
		connection, error := listener.Accept()

		if error != nil {
			test.Error("Accept error ", error)
		}

		connection_hybrid := Create_Conn(connection, server_tls_config)
		defer connection_hybrid.Close()
		buffer := make([]byte, 1024)
		length, error := connection_hybrid.Read(buffer)

		if error != nil {
			test.Error("Read error ", error)
		}

		test.Log("message", buffer[:length])

		if !connection_hybrid.Get_TLS() {
			test.Error("Get_TLS error")
		}

		if length != len(message) || !reflect.DeepEqual(message, buffer[:length]) {
			test.Error("messages not match ", message, buffer[:length])
		}

		break
	}
}
