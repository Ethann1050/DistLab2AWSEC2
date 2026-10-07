package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	for {
		response, _ := ln.Accept()
		conns <- response
	}
	// and add it to the channel for handling connections.
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	reader := bufio.NewReader(client)

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		msgs <- Message{
			sender:  clientid,
			message: text,
		}
	}
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	clientID := 0
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.
	ln, err := net.Listen("tcp", *portPtr)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return // Stop the program here so it doesn't crash later
	}
	fmt.Printf("Started Server")
	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			clientID += 1
			// - add the client to the clients map
			clients[clientID] = conn
			// - start to asynchronously handle messages from this client
			go handleClient(conn, clientID, msgs)
		case msg := <-msgs:
			//TODO Deal with a new message
			for clientID, connection := range clients {
				if msg.sender != clientID {
					fmt.Fprint(connection, msg.message)
				}

			}
			// Send the message to all clients that aren't the sender

		}

	}
}
