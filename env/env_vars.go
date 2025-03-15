package env

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/coopersmall/subswag/utils"
)

const (
	ENV_FILE = ".env"
)

type EnvVar string

const (
	API_URL     EnvVar = "API_URL"
	API_TIMEOUT EnvVar = "API_IMEOUT"

	APM_URL     EnvVar = "APM_URL"
	APM_TIMEOUT EnvVar = "APM_TIMEOUT"

	OPENAI_API_KEY EnvVar = "OPENAI_API_KEY"
	GROQ_API_KEY   EnvVar = "GROQ_API_KEY"

	POSTGRES_URL EnvVar = "POSTGRES_URL"

	REDIS_URL      EnvVar = "REDIS_URL"
	REDIS_PASSWORD EnvVar = "REDIS_PASSWORD"
	REDIS_DB       EnvVar = "REDIS_DB"

	REQUESTS_PUBLIC_KEY  EnvVar = "REQUESTS_PUBLIC_SIGNING_KEY"
	REQUESTS_PRIVATE_KEY EnvVar = "REQUESTS_PRIVATE_SIGNING_KEY"

	SECRETS_PUBLIC_KEY  EnvVar = "SECRETS_PUBLIC_SIGNING_KEY"
	SECRETS_PRIVATE_KEY EnvVar = "SECRETS_PRIVATE_SIGNING_KEY"

	JWT_SIGNING_KEY EnvVar = "JWT_SIGNING_KEY"
)

type IEnvVars interface {
	GetRequestsRSAPublicKey() (*rsa.PublicKey, error)
	GetRequestsRSAPrivateKey() (*rsa.PrivateKey, error)
	GetSecretsRSAPublicKey() (*rsa.PublicKey, error)
	GetSecretsRSAPrivateKey() (*rsa.PrivateKey, error)
	GetJWTSigningKey() ([]byte, error)
	GetRedisURL() (string, error)
	GetRedisPassword() (string, error)
	GetRedisDB() (int, error)
	GetAPMURL() (string, error)
	GetAPMTimeout() (time.Duration, error)
	GetSharedPostgresURL() (string, error)
	GetStandardPostgresURL() (string, error)
	GetOpenAIKey() (string, error)
	GetGroqKey() (string, error)
	GetAPIURL() (string, error)
	GetAPITimeout() (time.Duration, error)
}

type EnvVars struct {
	apmURL              *string
	apmTimeout          *time.Duration
	requestsPublicKey   *rsa.PublicKey
	requestsPrivateKey  *rsa.PrivateKey
	secretsPublicKey    *rsa.PublicKey
	secretsPrivateKey   *rsa.PrivateKey
	apiRouterURL        *string
	apiRouterTimeout    *time.Duration
	redisURL            *string
	redisPassword       *string
	redisDB             *int
	jwtsecret           []byte
	openaiAPIKey        *string
	groqAPIKey          *string
	sharedPostgresURL   *string
	standardPostgresURL *string
}

var Opts = []option{
	WithRedis(),
	WithAPM(),
	WithRequestsKeys(),
	WithSecretsKeys(),
	WithPostgres(),
	WithOpenAI(),
	WithGroq(),
	WithJWTSigner(),
	WithAPIConnParams(),
}

type option func(*EnvVars) error

func WithRedis() option {
	return func(e *EnvVars) error {
		url, err := GetEnvVar(REDIS_URL)
		if err != nil {
			return err
		}
		e.redisURL = &url
		return nil
	}
}

func WithAPM() option {
	return func(e *EnvVars) error {
		url, err := GetEnvVar(APM_URL)
		if err != nil {
			return err
		}
		timeout, err := GetEnvVar(APM_TIMEOUT)
		if err != nil {
			return err
		}
		parsedTimeout, err := constructDuration(timeout)

		e.apmURL = &url
		e.apmTimeout = &parsedTimeout
		return nil
	}
}

func WithRequestsKeys() option {
	return func(e *EnvVars) error {
		publicKey, err := GetEnvVar(REQUESTS_PUBLIC_KEY)
		if err != nil {
			return err
		}
		privateKey, err := GetEnvVar(REQUESTS_PRIVATE_KEY)
		if err != nil {
			return err
		}
		pub, err := constructPublicKey(publicKey)
		if err != nil {
			return err
		}
		priv, err := constructPrivateKey(privateKey)
		if err != nil {
			return err
		}
		e.requestsPublicKey = pub
		e.requestsPrivateKey = priv
		return nil
	}
}

func WithSecretsKeys() option {
	return func(e *EnvVars) error {
		publicKey, err := GetEnvVar(SECRETS_PUBLIC_KEY)
		if err != nil {
			return err
		}
		privateKey, err := GetEnvVar(SECRETS_PRIVATE_KEY)
		if err != nil {
			return err
		}
		pub, err := constructPublicKey(publicKey)
		if err != nil {
			return err
		}
		priv, err := constructPrivateKey(privateKey)
		if err != nil {
			return err
		}
		e.secretsPublicKey = pub
		e.secretsPrivateKey = priv
		return nil
	}
}

func WithPostgres() option {
	return func(e *EnvVars) error {
		shared, err := GetEnvVar(POSTGRES_URL)
		if err != nil {
			return err
		}
		standard, err := GetEnvVar(POSTGRES_URL)
		if err != nil {
			return err
		}
		e.sharedPostgresURL = &shared
		e.standardPostgresURL = &standard
		return nil
	}
}

func WithOpenAI() option {
	return func(e *EnvVars) error {
		key, err := GetEnvVar(OPENAI_API_KEY)
		if err != nil {
			return err
		}
		e.openaiAPIKey = &key
		return nil
	}
}

func WithGroq() option {
	return func(e *EnvVars) error {
		key, err := GetEnvVar(GROQ_API_KEY)
		if err != nil {
			return err
		}
		e.groqAPIKey = &key
		return nil
	}
}

func WithJWTSigner() option {
	return func(e *EnvVars) error {
		key, err := GetEnvVar(JWT_SIGNING_KEY)
		if err != nil {
			return err
		}
		e.jwtsecret = []byte(key)
		return nil
	}
}

func WithAPIConnParams() option {
	return func(e *EnvVars) error {
		apiRouterURL, err := GetEnvVar(API_URL)
		if err != nil {
			return err
		}
		apiRouterTimeout, err := GetEnvVar(API_TIMEOUT)
		if err != nil {
			return err
		}
		parsedTimeout, err := constructDuration(apiRouterTimeout)
		if err != nil {
			return err
		}
		e.apiRouterURL = &apiRouterURL
		e.apiRouterTimeout = &parsedTimeout
		return nil
	}
}

// Interface implementation methods
func (e *EnvVars) GetRequestsRSAPublicKey() (*rsa.PublicKey, error) {
	if e.requestsPublicKey == nil {
		return nil, errors.New("requests public key not initialized")
	}
	return e.requestsPublicKey, nil
}

func (e *EnvVars) GetRequestsRSAPrivateKey() (*rsa.PrivateKey, error) {
	if e.requestsPrivateKey == nil {
		return nil, errors.New("requests private key not initialized")
	}
	return e.requestsPrivateKey, nil
}

func (e *EnvVars) GetSecretsRSAPublicKey() (*rsa.PublicKey, error) {
	if e.secretsPublicKey == nil {
		return nil, errors.New("secrets public key not initialized")
	}
	return e.secretsPublicKey, nil
}

func (e *EnvVars) GetSecretsRSAPrivateKey() (*rsa.PrivateKey, error) {
	if e.secretsPrivateKey == nil {
		return nil, errors.New("secrets private key not initialized")
	}
	return e.secretsPrivateKey, nil
}

func (e *EnvVars) GetJWTSigningKey() ([]byte, error) {
	if e.jwtsecret == nil {
		return nil, errors.New("JWT signing key not initialized")
	}
	return e.jwtsecret, nil
}

func (e *EnvVars) GetRedisURL() (string, error) {
	if e.redisURL == nil {
		return "", errors.New("redis URL not initialized")
	}
	return *e.redisURL, nil
}

func (e *EnvVars) GetRedisPassword() (string, error) {
	if e.redisPassword == nil {
		return "", errors.New("redis password not initialized")
	}
	return *e.redisPassword, nil
}

func (e *EnvVars) GetRedisDB() (int, error) {
	if e.redisDB == nil {
		return 0, errors.New("redis DB not initialized")
	}
	return *e.redisDB, nil
}

func (e *EnvVars) GetAPMURL() (string, error) {
	if e.apmURL == nil {
		return "", errors.New("APM URL not initialized")
	}
	return *e.apmURL, nil
}

func (e *EnvVars) GetAPMTimeout() (time.Duration, error) {
	if e.apmTimeout == nil {
		return 0, errors.New("APM timeout not initialized")
	}
	return *e.apmTimeout, nil
}

func (e *EnvVars) GetSharedPostgresURL() (string, error) {
	if e.sharedPostgresURL == nil {
		return "", errors.New("shared postgres URL not initialized")
	}
	return *e.sharedPostgresURL, nil
}

func (e *EnvVars) GetStandardPostgresURL() (string, error) {
	if e.standardPostgresURL == nil {
		return "", errors.New("standard postgres URL not initialized")
	}
	return *e.standardPostgresURL, nil
}

func (e *EnvVars) GetOpenAIKey() (string, error) {
	if e.openaiAPIKey == nil {
		return "", errors.New("OpenAI API key not initialized")
	}
	return *e.openaiAPIKey, nil
}

func (e *EnvVars) GetGroqKey() (string, error) {
	if e.groqAPIKey == nil {
		return "", errors.New("Groq API key not initialized")
	}
	return *e.groqAPIKey, nil
}

func (e *EnvVars) GetAPIURL() (string, error) {
	if e.apiRouterURL == nil {
		return "", errors.New("API router URL not initialized")
	}
	return *e.apiRouterURL, nil
}

func (e *EnvVars) GetAPITimeout() (time.Duration, error) {
	if e.apiRouterTimeout == nil {
		return 0, errors.New("API router timeout not initialized")
	}
	return *e.apiRouterTimeout, nil
}

// Constructor functions
func GetEnvVars(opts ...option) (IEnvVars, error) {
	vars := &EnvVars{}
	for _, opt := range opts {
		if err := opt(vars); err != nil {
			return nil, err
		}
	}
	return vars, nil
}

func MustGetEnvVars(opts ...option) IEnvVars {
	vars, err := GetEnvVars(opts...)
	if err != nil {
		panic(err)
	}
	return vars
}

func GetEnvVar(key EnvVar) (string, error) {
	found, ok := os.LookupEnv(string(key))
	if !ok {
		return "", utils.NewInvalidArgumentError("environment variable not found")
	}
	return found, nil
}

func constructPublicKey(key string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("failed to decode public key")
	}
	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		return key.(*rsa.PublicKey), nil
	}
	return nil, errors.New("failed to parse public key")
}

func constructPrivateKey(raw string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, errors.New("failed to decode private key")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, err
	}

	return rsaKey, nil
}

func constructDuration(raw string) (time.Duration, error) {
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	return time.Duration(parsed), nil
}
