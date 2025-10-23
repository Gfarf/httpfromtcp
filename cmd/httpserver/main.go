package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Gfarf/httpfromtcp/internal/headers"
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

func writingHeaders(w *response.Writer, req *request.Request) {
	var sCode response.StatusCode
	var buf []byte
	h := response.GetDefaultHeaders(0)
	if req.RequestLine.RequestTarget == "/yourproblem" {
		buf = []byte("<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>")
		sCode = 400
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		buf = []byte("<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>")
		sCode = 500
	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		handlerHTTPBIN(w, req, h)
		return
	} else if req.RequestLine.RequestTarget == "/video" {
		handlerVideo(w, req, h)
		return
	} else {
		buf = []byte("<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>")
		sCode = 200
	}

	h.Override("content-length", fmt.Sprintf("%d", len(buf)))
	h.Override("content-type", "text/html")
	w.WriteStatusLine(sCode)
	w.WriteHeaders(h)
	w.WriteBody(buf)
}

func handlerHTTPBIN(w *response.Writer, req *request.Request, h headers.Headers) {
	targetURL := "https://httpbin.org/" + strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	delete(h, "content-length")
	h.Override("Transfer-Encoding", "chunked")
	h.Override("Trailer", "X-Content-Sha256, X-Content-Length")
	resp, err := http.Get(targetURL)
	if err != nil {
		w.WriteStatusLine(400)
		body := []byte(fmt.Sprintf("Error connecting to httpbin: %v\n", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	defer resp.Body.Close()
	w.WriteStatusLine(200)
	w.WriteHeaders(h)
	buf := make([]byte, bufferSize)
	full := make([]byte, 0)
	for {
		n, err := resp.Body.Read(buf)
		fmt.Printf("%d bytes read\n", n)
		if n > 0 {
			w.WriteChunkedBody(buf[:n])
			full = append(full, buf[:n]...)
		}
		if err == io.EOF {
			fmt.Printf("%d bytes read\n", n)
			w.WriteChunkedBodyDone()
			break
		}
		if err != nil {
			body := []byte(fmt.Sprintf("Error reading from body: %v\n", err))
			w.WriteBody(body)
			return
		}
	}
	trailers := headers.NewHeaders()
	sha256 := fmt.Sprintf("%x", sha256.Sum256(full))
	trailers.OverrideTrailers("X-Content-Sha256", sha256)
	trailers.OverrideTrailers("X-Content-Length", fmt.Sprintf("%d", len(full)))
	err = w.WriteTrailers(trailers)
	if err != nil {
		fmt.Println("Error writing trailers:", err)
	}
	fmt.Println("Wrote trailers")
}

func handlerVideo(w *response.Writer, req *request.Request, h headers.Headers) {
	h.Override("content-type", "video/mp4")
	f, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		w.WriteStatusLine(500)
		body := []byte(fmt.Sprintf("Error reading file: %v\n", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	h.Override("content-length", fmt.Sprintf("%d", len(f)))
	w.WriteStatusLine(200)
	w.WriteHeaders(h)
	w.WriteBody(f)
}
