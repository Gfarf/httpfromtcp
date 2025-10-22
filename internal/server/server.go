package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync/atomic"

	"github.com/Gfarf/httpfromtcp/internal/request"
	"github.com/Gfarf/httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

// Contains the state of the server
type Server struct {
	closed      atomic.Bool
	listener    net.Listener
	handlerFunc Handler
}

const bufferSize = 1024

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
	lineChannels, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	herr := s.handlerFunc(conn, lineChannels)
	if err != nil {
		writeErrorHandler(conn, herr)
	}
}

func writeErrorHandler(w io.Writer, err *HandlerError) error {
	output := fmt.Sprintf("error %s, code %v", err.Message, err.StatusCode)
	_, err1 := w.Write([]byte(output))
	if err1 != nil {
		return err1
	}
	return nil
}
