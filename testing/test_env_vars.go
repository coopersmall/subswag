package testing

import (
	"crypto/rand"
	"crypto/rsa"
	"time"
)

type testEnvVars struct {
	requestsPublicKey   *rsa.PublicKey
	requestsPrivateKey  *rsa.PrivateKey
	secretsPublicKey    *rsa.PublicKey
	secretsPrivateKey   *rsa.PrivateKey
	jwtSigningKey       []byte
	redisURL            string
	redisPassword       string
	redisDB             int
	apmURL              string
	apmTimeout          time.Duration
	sharedPostgresURL   string
	standardPostgresURL string
	openAIKey           string
	groqKey             string
	apiURL              string
	apiTimeout          time.Duration
}

func NewTestEnvVars() (*testEnvVars, error) {
	requestsPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	requestsPublicKey := &requestsPrivateKey.PublicKey

	secretsPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	secretsPublicKey := &secretsPrivateKey.PublicKey

	return &testEnvVars{
		requestsPublicKey:   requestsPublicKey,
		requestsPrivateKey:  requestsPrivateKey,
		secretsPublicKey:    secretsPublicKey,
		secretsPrivateKey:   secretsPrivateKey,
		jwtSigningKey:       []byte("test-jwt-signing-key"),
		redisURL:            "localhost:6379",
		redisPassword:       "test-password",
		redisDB:             0,
		apmURL:              "http://localhost:8200",
		apmTimeout:          time.Second * 5,
		sharedPostgresURL:   "postgres://test:test@localhost:5432/testdb?sslmode=disable",
		standardPostgresURL: "postgres://test:test@localhost:5432/testdb?sslmode=disable",
		openAIKey:           "test-openai-key",
		groqKey:             "test-groq-key",
		apiURL:              "http://localhost:8080",
		apiTimeout:          time.Second * 10,
	}, nil
}

// Implementation of IEnvVars interface
func (t *testEnvVars) GetRequestsRSAPublicKey() (*rsa.PublicKey, error) {
	return t.requestsPublicKey, nil
}

func (t *testEnvVars) GetRequestsRSAPrivateKey() (*rsa.PrivateKey, error) {
	return t.requestsPrivateKey, nil
}

func (t *testEnvVars) GetSecretsRSAPublicKey() (*rsa.PublicKey, error) {
	return t.secretsPublicKey, nil
}

func (t *testEnvVars) GetSecretsRSAPrivateKey() (*rsa.PrivateKey, error) {
	return t.secretsPrivateKey, nil
}

func (t *testEnvVars) GetJWTSigningKey() ([]byte, error) {
	return t.jwtSigningKey, nil
}

func (t *testEnvVars) GetRedisURL() (string, error) {
	return t.redisURL, nil
}

func (t *testEnvVars) GetRedisPassword() (string, error) {
	return t.redisPassword, nil
}

func (t *testEnvVars) GetRedisDB() (int, error) {
	return t.redisDB, nil
}

func (t *testEnvVars) GetAPMURL() (string, error) {
	return t.apmURL, nil
}

func (t *testEnvVars) GetAPMTimeout() (time.Duration, error) {
	return t.apmTimeout, nil
}

func (t *testEnvVars) GetSharedPostgresURL() (string, error) {
	return t.sharedPostgresURL, nil
}

func (t *testEnvVars) GetStandardPostgresURL() (string, error) {
	return t.standardPostgresURL, nil
}

func (t *testEnvVars) GetOpenAIKey() (string, error) {
	return t.openAIKey, nil
}

func (t *testEnvVars) GetGroqKey() (string, error) {
	return t.groqKey, nil
}

func (t *testEnvVars) GetAPIURL() (string, error) {
	return t.apiURL, nil
}

func (t *testEnvVars) GetAPITimeout() (time.Duration, error) {
	return t.apiTimeout, nil
}

// Helper methods to update test environment variables
func (t *testEnvVars) SetRedisURL(url string) {
	t.redisURL = url
}

func (t *testEnvVars) SetRedisPassword(password string) {
	t.redisPassword = password
}

func (t *testEnvVars) SetSharedPostgresURL(url string) {
	t.sharedPostgresURL = url
}

func (t *testEnvVars) SetStandardPostgresURL(url string) {
	t.standardPostgresURL = url
}

func (t *testEnvVars) SetAPMURL(url string) {
	t.apmURL = url
}
