package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/quic-go/quic-go"
)

const (
	clientAddr = "127.0.0.1:7000"
	serverAddr = "127.0.0.1:1999"
)

var (
	mu     sync.Mutex
	quicCh = make(chan string) // Channel for sending QUIC messages
	httpCh = make(chan string) // Channel for sending HTTP response
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Request body is a JSON string
	fmt.Println("Received request with JSON body: ", string(body))

	// Send to the QUIC server
	quicCh <- string(body)

	// Response from the QUIC server
	res := <-httpCh
	var jsonData = map[string]string{
		"type":    "FRAME",
		"message": res,
	}
	jsonRes, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}
	w.Write(jsonRes)

	// Respond with a JSON object
	//w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte(`{"type": "FRAME","message": "from clientServer"}`))
	//w.Write([]byte(res))
}
func main() {
	// API SETUP
	go func() {
		http.HandleFunc("/", handleFunc)
		http.ListenAndServe(clientAddr, nil)
	}()
	fmt.Println("HTTP Server listening on", clientAddr)

	// CONNECT TO QUIC SERVER
	// QUIC Config
	quicConf := &quic.Config{
		EnableDatagrams: true, // 0-RTT
	}

	// Connect to the server
	conn, err := quic.DialAddr(context.Background(), serverAddr, &tls.Config{InsecureSkipVerify: true}, quicConf)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to server.")
	//fmt.Print("Enter message: ")

	go func() {
		for {
			stream, err := conn.AcceptStream(context.Background())
			if err != nil {
				panic(err)
			}
			go handleStream(stream)
		}
	}()

	for { // keep sending the message from the web
		//reader := bufio.NewReader(os.Stdin)
		//line, _ := reader.ReadString('\n')
		line := <-quicCh
		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(stream, line)
		fmt.Println("SENT: ", line)
	}
}

func handleStream(stream quic.Stream) { // keep sending back to the web
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

		fmt.Println("Received:", string(data[:n]))
		httpCh <- string(data[:n])
	}
}
