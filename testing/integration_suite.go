package testing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IntegrationSuiteConfig struct {
	PostgresURL   string
	RedisURL      string
	RedisPassword string
}

type IntegrationTest struct {
	suite.Suite
	*TestEnv
	config IntegrationSuiteConfig
}

// New approach: each test creates its own environment
func NewIntegrationTest(
	t *testing.T,
	config IntegrationSuiteConfig,
) *IntegrationTest {
	test := new(IntegrationTest)
	test.config = config
	test.SetT(t)
	return test
}

func (s *IntegrationTest) BeforeTest(suiteName, testName string) {
	if err := s.Reset(); err != nil {
		s.T().Fatalf("Failed to reset test environment: %v", err)
	}
}

func (s *IntegrationTest) SetupTest() {
	env, err := s.createEnv()
	if err != nil {
		s.T().Fatalf("Failed to create test environment: %v", err)
	}
	if err := env.start(); err != nil {
		s.T().Fatalf("Failed to start test environment: %v", err)
	}
	s.TestEnv = env
}

func (s *IntegrationTest) TearDownTest() {
	if s.TestEnv != nil {
		if err := s.TestEnv.stop(); err != nil {
			s.T().Errorf("Failed to stop test environment: %v", err)
		}
	}
	s.TestEnv = nil
}

func (s *IntegrationTest) createEnv() (*TestEnv, error) {
	opts := []TestEnvOption{}
	if s.config.PostgresURL != "" {
		opts = append(opts, WithExternalPostgres(s.config.PostgresURL))
	}
	if s.config.RedisURL != "" {
		opts = append(opts, WithExternalRedis(s.config.RedisURL, s.config.RedisPassword))
	}

	return NewTestingEnvWithOptions(opts...)
}

func RunIntegrationTest[T suite.TestingSuite](
	t *testing.T,
	config IntegrationSuiteConfig,
	factory func(*IntegrationTest) T,
) {
	t.Parallel()
	baseTest := NewIntegrationTest(t, config)
	suite.Run(t, factory(baseTest))
}

func GetIntegrationSuiteConfig() IntegrationSuiteConfig {
	if os.Getenv("CI") == "true" {
		return getCIIntegrationTestConfig()
	}
	return getDefaultIntegrationSuiteConfig()
}

func getDefaultIntegrationSuiteConfig() IntegrationSuiteConfig {
	return IntegrationSuiteConfig{
		PostgresURL:   "postgres://test:test@localhost:5432/testdb?sslmode=disable",
		RedisURL:      "localhost:6379",
		RedisPassword: "",
	}
}

func getCIIntegrationTestConfig() IntegrationSuiteConfig {
	return IntegrationSuiteConfig{
		PostgresURL:   os.Getenv("CI_POSTGRES_URL"),
		RedisURL:      os.Getenv("CI_REDIS_URL"),
		RedisPassword: os.Getenv("CI_REDIS_PASSWORD"),
	}
}
