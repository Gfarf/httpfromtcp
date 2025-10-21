package response

import (
	"fmt"
	"io"

	"github.com/Gfarf/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	statusCodeOK StatusCode = iota
	statusCodeInternalServerError
	statusCodeBadRequest
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case statusCodeOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
	case statusCodeBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
	case statusCodeInternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
	default:
		_, err := w.Write([]byte(""))
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = fmt.Sprintf("%d", contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		fullText := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Write([]byte(fullText))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}
