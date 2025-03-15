package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *HttpClientTestSuite) TestGet() {
	s.Run("successful GET request", func() {
		ctx := context.Background()
		result, err := s.client.Get(ctx, "/test", nil)

		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)

		responseMap, ok := result.(map[string]interface{})
		assert.True(s.T(), ok)
		assert.Equal(s.T(), "test response", responseMap["message"])
	})

	s.Run("GET request with error response", func() {
		ctx := context.Background()
		result, err := s.client.Get(ctx, "/error", nil)

		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)
		assert.Contains(s.T(), err.Error(), "unexpected status code: 500")
	})

	s.Run("GET request with context cancellation", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel the context immediately

		result, err := s.client.Get(ctx, "/test", nil)

		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)
		assert.Contains(s.T(), err.Error(), "context canceled")
	})
}

func (s *HttpClientTestSuite) TestPost() {
	s.Run("successful POST request", func() {
		ctx := context.Background()
		body := map[string]string{"key": "value"}
		result, err := s.client.Post(ctx, "/test", body, nil)

		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)

		responseMap, ok := result.(map[string]interface{})
		assert.True(s.T(), ok)
		assert.Equal(s.T(), "test response", responseMap["message"])
	})
}

func (s *HttpClientTestSuite) TestFetchRawResponse() {
	s.Run("successful raw response fetch", func() {
		ctx := context.Background()
		resp, err := s.client.FetchRawResponse(ctx, "/test", http.MethodGet, nil, nil)

		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), resp)
		assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

		var result map[string]string
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "test response", result["message"])
	})
}

func (s *HttpClientTestSuite) TestRequestWithHeaders() {
	s.Run("request with custom headers", func() {
		ctx := context.Background()
		headers := map[string]string{
			"X-Custom-Header": "custom value",
		}

		_, err := s.client.Get(ctx, "/test", headers)

		assert.NoError(s.T(), err)
		// In a real scenario, you might want to check if the server received the custom header
	})
}

func (s *HttpClientTestSuite) TestRequestWithTimeout() {
	s.Run("request with timeout", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		// Assuming the server is configured to sleep for longer than the timeout
		_, err := s.client.Get(ctx, "/slow", nil)

		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "context deadline exceeded")
	})
}
