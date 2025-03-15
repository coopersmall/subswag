package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpclient "github.com/coopersmall/subswag/clients/http"
	"github.com/stretchr/testify/suite"
)

type HttpClientTestSuite struct {
	suite.Suite
	server *httptest.Server
	client *httpclient.HttpClient
}

func (s *HttpClientTestSuite) SetupTest() {
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/test":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "test response"})
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
		case "/slow":
			time.Sleep(100 * time.Millisecond)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	s.client = httpclient.NewHttpClient(s.server.URL, nil)
}

func (s *HttpClientTestSuite) TearDownTest() {
	s.server.Close()
}

func TestHttpClientSuite(t *testing.T) {
	suite.Run(t, new(HttpClientTestSuite))
}
