package server_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/http/server"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()

	userId        = user.UserID(12345)
	tokenId       = apitoken.APITokenID(12345)
	validUserData = user.UserData{
		Email: "joe@joe.com",
	}
	validTokenData = &apitoken.APITokenData{
		UserId:      userId,
		Expiry:      time.Now().Add(time.Hour),
		Permissions: []domain.Permission{domain.APIPermission},
	}
	routeFunc = func(req server.IRequest) (any, error) {
		return userId, nil
	}
	validToken *apitoken.APITokenWithSecret
)

func (s *HTTPServerTestSuite) SetupSubTest() {
	s.Reset()
	services, _ := s.GetServices()
	_, err := services.UsersService().CreateUserWithId(ctx, userId, validUserData)
	assert.NoError(s.T(), err)
	token, err := services.APITokenService(userId).CreateTokenWithId(ctx, utils.ID(tokenId), validTokenData)
	assert.NoError(s.T(), err)
	validToken = token
}

func (s *HTTPServerTestSuite) TestAPIRouteSuccess() {
	s.Run("works with authorization header", func() {
		handler := server.NewAPIRoute(
			"/test",
			"",
			http.MethodGet,
			routeFunc,
			s.TestEnv,
		)

		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		assert.NoError(s.T(), err)
		req.Header.Set("Authorization", "Bearer "+validToken.Secret)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(s.T(), http.StatusOK, rr.Code)
		body := rr.Body.String()
		assert.Equal(s.T(), body, fmt.Sprintf("%d", userId))
	})

	s.Run("returns error for invalid token", func() {
		handler := server.NewAPIRoute(
			"/test",
			"",
			http.MethodGet,
			routeFunc,
			s.TestEnv,
		)

		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		assert.NoError(s.T(), err)
		req.Header.Set("Authorization", "Bearer invalid_token")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(s.T(), http.StatusUnauthorized, rr.Code)
		assert.Equal(s.T(), "null", rr.Body.String())
	})

	s.Run("returns error when route function fails", func() {
		handler := server.NewAPIRoute(
			"/test",
			"",
			http.MethodGet,
			func(req server.IRequest) (any, error) {
				return nil, fmt.Errorf("route function error")
			},
			s.TestEnv,
		)

		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		assert.NoError(s.T(), err)
		req.Header.Set("Authorization", "Bearer "+validToken.Secret)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(s.T(), http.StatusInternalServerError, rr.Code)
		assert.Equal(s.T(), "Internal Server Error", rr.Body.String())
	})

	s.Run("returns error when token retrieval fails", func() {
		services, _ := s.GetServices()
		services.APITokenService(userId).DeleteToken(ctx, validToken.ID)

		handler := server.NewAPIRoute(
			"/test",
			"",
			http.MethodGet,
			routeFunc,
			s.TestEnv,
		)

		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		assert.NoError(s.T(), err)
		req.Header.Set("Authorization", "Bearer "+validToken.Secret)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(s.T(), http.StatusInternalServerError, rr.Code)
		assert.Equal(s.T(), "Internal Server Error", rr.Body.String())
	})
}
