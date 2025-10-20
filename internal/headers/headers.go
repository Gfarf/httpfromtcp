package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}
	idxColon := bytes.Index(data, []byte(":"))
	if idxColon > 0 {
		if data[idxColon-1] == ' ' {
			return 0, false, fmt.Errorf("poorly formatted field name: %s", string(data))
		}
	} else {
		return 0, false, fmt.Errorf("poorly formatted field line: %s", string(data))
	}
	fieldName := data[:idxColon]
	fieldValue := data[idxColon+1 : idx]
	value := strings.TrimSpace(string(fieldValue))
	value = strings.TrimRight(value, "\r\n")
	key := strings.ToLower(strings.TrimSpace(string(fieldName)))
	if !validTokens([]byte(key)) {
		return 0, false, fmt.Errorf("invalid header token found: %s", key)
	}
	h.Set(key, string(value))
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	oldValue, ok := h[key]
	if ok {
		h[key] = oldValue + ", " + value
	} else {
		h[key] = value
	}
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

// validTokens checks if the data contains only valid tokens
// or characters that are allowed in a token
func validTokens(data []byte) bool {
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' {
		return true
	}

	return slices.Contains(tokenChars, c)
}

func (h Headers) Get(key string) string {
	value, ok := h[strings.ToLower(key)]
	if ok {
		return value
	}
	return ""
}
