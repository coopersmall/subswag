package server_test

import (
	"context"
	"net/http"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *HTTPServerTestSuite) TestGetCorrelationIdFromRequest() {
	s.Run("CorrelationID from context", func() {
		expectedID := domain.NewCorrelationID()
		ctx := context.WithValue(context.Background(), server.CorrelationIDKey, expectedID)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		require.NoError(s.T(), err)
		result := server.GetCorrelationIdFromRequest(req)
		assert.Equal(s.T(), expectedID, result)
	})

	s.Run("CorrelationID from header", func() {
		expectedID := domain.NewCorrelationID()
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(s.T(), err)
		req.Header.Set(server.CorrelationIDHeader, expectedID.String())
		result := server.GetCorrelationIdFromRequest(req)
		assert.Equal(s.T(), expectedID, result)
	})

	s.Run("New CorrelationID when not in context or header", func() {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(s.T(), err)
		result := server.GetCorrelationIdFromRequest(req)
		assert.NotEmpty(s.T(), result)
	})

	s.Run("Invalid CorrelationID in header", func() {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(s.T(), err)
		req.Header.Set(server.CorrelationIDHeader, "invalid-correlation-id")
		result := server.GetCorrelationIdFromRequest(req)
		assert.NotEmpty(s.T(), result)
		assert.NotEqual(s.T(), "invalid-correlation-id", result.String())
	})
}
