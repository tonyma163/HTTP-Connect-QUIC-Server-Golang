package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"sync"

	quic "github.com/quic-go/quic-go"
)

const addr = "127.0.0.1:1999"

var (
	mu          sync.Mutex
	connections = make(map[quic.Connection]struct{}) // Saving all connections for the connected clients
)

func main() {
	// QUIC Config
	quicConf := &quic.Config{
		EnableDatagrams: true, // 0-RTT
	}

	listener, err := quic.ListenAddr(addr, generateTLSConfig(), quicConf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server listening on %s\n", addr)

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			panic(err)
		}

		// Add the connection into the map
		mu.Lock()
		connections[conn] = struct{}{}
		mu.Unlock()

		go handleConnection(conn)
	}
}

func handleConnection(conn quic.Connection) {
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			panic(err)
		}

		go func() {
			defer stream.Close()

			r := bufio.NewReader(stream)
			for {
				//line, err := r.ReadString('\n')
				data := make([]byte, 1024)
				n, err := r.Read(data)
				if err == io.EOF {
					break
				} else if err != nil {
					panic(err)
				}

				var jsonData map[string]string
				err = json.Unmarshal(data[:n], &jsonData)

				jsonMsg, err := json.Marshal(jsonData)
				if err != nil {
					panic(err)
				}

				fmt.Print(jsonMsg)
				// Broadcast the message to all other connected clients
				broadcast(string(jsonMsg), conn)
			}

			// Remove connection when the client disconnects
			mu.Lock()
			delete(connections, conn)
			mu.Unlock()
		}()
	}
}

func broadcast(message string, sender quic.Connection) {
	mu.Lock()
	defer mu.Unlock()

	//
	//message = "From: " + sender.RemoteAddr().String() + ": " + message

	//
	for conn := range connections {
		if conn != sender { // Skip the sender
			stream, err := conn.OpenStream()
			if err != nil {
				fmt.Println("Could not open stream to client:", err)
				continue
			}
			fmt.Fprintf(stream, message)
			stream.Close()
			fmt.Println("SENT")
		}
	}
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}, InsecureSkipVerify: true}
}
