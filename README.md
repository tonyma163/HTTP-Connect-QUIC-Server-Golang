# HTTP-Connect-QUIC-Server-Golang

The web side (javascript file) will keep sending/receiving the json from the client side server using api fetching <br />
while the QUIC server will keep sending the received json from the client side server to other clients (except the sender client) in the same room.

## Simple Data Flow
ClientA <---> [Web1.js] <---JSON---> [Client1 HTTP Server] <br />
<---JSON_STRING---> [QUIC SERVER] <---JSON_STRING---> <br />
[Client2 HTTP Server] <---JSON---> [Web2.js] <---> ClientB <br />

# Steps
## 1 Run QUIC Server
go run server.go

## 2 Run the Both Client Server
```
cd client
go run client.go
```
cd client2
go run client2.go

## 3 Run the Both Javascript File (Send Json & Read)
cd client
node web.js <br />
cd client2 <br />
node web3.js

# PORT
server.go 1999 <br />
client.go 6969 <br />
client2.go 7000
