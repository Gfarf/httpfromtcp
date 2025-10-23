package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/Gfarf/httpfromtcp/internal/request"
	"github.com/Gfarf/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

// Contains the state of the server
type Server struct {
	closed      atomic.Bool
	listener    net.Listener
	handlerFunc Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	//Creates a net.Listener and returns a new Server instance. Starts listening for requests inside a goroutine.
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := Server{
		listener:    l,
		handlerFunc: handler,
	}
	go s.listen()
	return &s, nil
}

func (s *Server) Close() error {
	//Closes the listener and the server
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	//Uses a loop to .Accept new connections as they come in, and handles each one in a new goroutine.
	//I used an atomic.Bool to track whether the server is closed or not so that I can ignore connection errors after the server is closed.
	for {
		// Wait for a connection.
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	w := response.NewWriter(conn)
	lineChannels, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(400)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handlerFunc(w, lineChannels)
}
