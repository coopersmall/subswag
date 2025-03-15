package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/coopersmall/subswag/utils"
)

type server struct {
	logger       utils.ILogger
	readTimeout  time.Duration
	writeTimeout time.Duration
	router       http.Handler
}

func (s *server) Serve(conn net.Conn) {
	now := time.Now()
	defer conn.Close()

	conn.SetReadDeadline(now.Add(s.readTimeout))
	conn.SetWriteDeadline(now.Add(s.writeTimeout))

	r, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return
	}

	w := newWriter(conn)
	s.router.ServeHTTP(w, r)

	statusCode := w.StatusCode()
	if statusCode >= 500 {
		statusText := strings.ToLower(http.StatusText(statusCode))
		s.logger.Error(
			r.Context(),
			"failed request",
			errors.New(fmt.Sprintf("failed request: %s", http.StatusText(statusCode))),
			map[string]any{
				"status": statusCode,
				"text":   statusText,
			},
		)
	}
}

func newServer(
	router *Router,
	readTimeout time.Duration,
	writeTimeout time.Duration,
) *server {
	return &server{
		logger:       utils.GetLogger("http-server"),
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		router:       router,
	}
}
