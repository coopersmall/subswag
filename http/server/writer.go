package server

import (
	"fmt"
	"net"
	"net/http"
)

type writer struct {
	header     http.Header
	conn       net.Conn
	statusCode int
}

func newWriter(conn net.Conn) *writer {
	return &writer{
		conn:   conn,
		header: make(http.Header),
	}
}

func (w *writer) configure(conn net.Conn) {
	w.conn = conn
}

func (w *writer) Header() http.Header {
	return w.header
}

func (w *writer) WriteHeader(statusCode int) {
	statusText := http.StatusText(statusCode)
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)

	w.statusCode = statusCode

	w.conn.Write([]byte(statusLine))

	for key, values := range w.header {
		for _, value := range values {
			headerLine := fmt.Sprintf("%s: %s\r\n", key, value)
			w.conn.Write([]byte(headerLine))
		}
	}

	w.conn.Write([]byte("\r\n"))
}

func (w *writer) Write(b []byte) (int, error) {
	return w.conn.Write(b)
}

func (w *writer) StatusCode() int {
	return w.statusCode
}
