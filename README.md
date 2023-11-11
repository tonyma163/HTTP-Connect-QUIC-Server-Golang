# HTTP-Connect-QUIC-Server-Golang

The web side (javascript file) will keep sending/receiving the json from the client side server using api fetching
While the QUIC server will keep sending the received json from the client side server to other clients (except the sender client) in the same room.

# Simple Data Flow
ClientA <---> Web1.js <---JSON---> Client1 HTTP Server <---JSON_STRING---> QUIC SERVER <---JSON_STRING---> Client2 HTTP Server <---JSON---> Web2.js <---> ClientB

# Steps
# 1 Run QUIC Server
go run server.go

# 2 Run the Both Client Server
cd client
go run client.go
cd client2
go run client2.go

# 3 Run the Both Javascript File (Send Json & Read)
cd client
node web.js
cd client2
node web3.js