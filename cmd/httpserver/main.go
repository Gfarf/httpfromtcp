package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gfarf/httpfromtcp/internal/request"
	"github.com/Gfarf/httpfromtcp/internal/response"
	"github.com/Gfarf/httpfromtcp/internal/server"
)

const port = 42069
const bufferSize = 1024

func main() {
	server, err := server.Serve(port, writingHeaders)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func writingHeaders(w io.Writer, req *request.Request) *server.HandlerError {
	errHand := server.HandlerError{}
	var sCode response.StatusCode
	buf := make([]byte, bufferSize)
	if req.RequestLine.RequestTarget == "/yourproblem" {
		buf = []byte("Your problem is not my problem\n")
		sCode = 400
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		buf = []byte("Woopsie, my bad\n")
		sCode = 500
	} else {
		buf = []byte("All good, frfr\n")
		sCode = 200
	}
	err := response.WriteStatusLine(w, sCode)
	if err != nil {
		errHand.Message = err.Error()
		errHand.StatusCode = 505
		return &errHand
	}
	err = response.WriteHeaders(w, response.GetDefaultHeaders(len(buf)))
	if err != nil {
		errHand.Message = err.Error()
		errHand.StatusCode = 505
		return &errHand
	}
	_, err = w.Write(buf)
	if err != nil {
		errHand.Message = err.Error()
		errHand.StatusCode = 505
		return &errHand
	}
	return nil
}
