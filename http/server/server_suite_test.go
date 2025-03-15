package server_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type HTTPServerTestSuite struct {
	*tt.IntegrationTest
}

func TestHTTPServerTestSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *HTTPServerTestSuite {
		return &HTTPServerTestSuite{
			IntegrationTest: its,
		}
	})
}
