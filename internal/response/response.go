package response

import (
	"fmt"
	"io"

	"github.com/Gfarf/httpfromtcp/internal/headers"
)

type Writer struct {
	statusWrite WriteCode
	Writer      io.Writer
}

type WriteCode int

const (
	statusWriteStatusLine WriteCode = iota
	statusWriteHeaders
	statusWriteBody
	statusWriteDone
)

type StatusCode int

const (
	statusCodeOK                  StatusCode = 200
	statusCodeInternalServerError StatusCode = 500
	statusCodeBadRequest          StatusCode = 400
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = fmt.Sprintf("%d", contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{
		statusWrite: statusWriteStatusLine,
		Writer:      conn,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.statusWrite != statusWriteStatusLine {
		return fmt.Errorf("must start writing from Status Line, then Headers, then Body")
	}
	switch statusCode {
	case statusCodeOK:
		_, err := w.Writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
	case statusCodeBadRequest:
		_, err := w.Writer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
	case statusCodeInternalServerError:
		_, err := w.Writer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
	default:
		_, err := w.Writer.Write([]byte(""))
		if err != nil {
			return err
		}
	}
	w.statusWrite = statusWriteHeaders
	return nil
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.statusWrite != statusWriteHeaders {
		return fmt.Errorf("must start writing from Status Line, then Headers, then Body")
	}
	for key, value := range headers {
		fullText := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Writer.Write([]byte(fullText))
		if err != nil {
			return err
		}
	}
	_, err := w.Writer.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	w.statusWrite = statusWriteBody
	return nil
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.statusWrite != statusWriteBody {
		return 0, fmt.Errorf("must start writing from Status Line, then Headers, then Body")
	}
	n, err := w.Writer.Write(p)
	if err != nil {
		return 0, err
	}
	w.statusWrite = statusWriteDone
	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.statusWrite != statusWriteBody {
		return 0, fmt.Errorf("must start writing from Status Line, then Headers, then Body")
	}
	n1, err := fmt.Fprintf(w.Writer, "%X\r\n", len(p))
	if err != nil {
		return n1, err
	}
	n2, err := w.Writer.Write(p)
	if err != nil {
		return n2, err
	}
	n3, err := w.Writer.Write([]byte("\r\n"))
	n := n1 + n2 + n3
	return n, err
}
func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.Writer.Write([]byte("0\r\n"))
	w.statusWrite = statusWriteDone
	return n, err
}
func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.statusWrite != statusWriteDone {
		return fmt.Errorf("must start writing from Status Line, then Headers, then Body, Then Trailers")
	}
	defer func() { w.statusWrite = statusWriteBody }()
	for k, v := range h {
		_, err := w.Writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.Writer.Write([]byte("\r\n"))
	return err
}
