package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/Gfarf/httpfromtcp/internal/request"
	"github.com/Gfarf/httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

// Contains the state of the server
type Server struct {
	closed   atomic.Bool
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	//Creates a net.Listener and returns a new Server instance. Starts listening for requests inside a goroutine.
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := Server{
		listener: l,
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
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	fmt.Println("handling something, starting to write status line")
	response.WriteStatusLine(conn, 0)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	return
}
