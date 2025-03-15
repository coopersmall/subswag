package utils

import (
	"errors"
	"os"
)

const (
	ENV_FILE = ".env"
)

type EnvVar string

func (e EnvVar) String() string {
	return string(e)
}

const (
	POSTGRES_HOST     EnvVar = "POSTGRES_HOST"
	POSTGRES_PORT     EnvVar = "POSTGRES_PORT"
	POSTGRES_USER     EnvVar = "POSTGRES_USER"
	POSTGRES_PASSWORD EnvVar = "POSTGRES_PASSWORD"
	POSTGRESS_SSL     EnvVar = "POSTGRES_SSL_MODE"

	HTTP_PORT          EnvVar = "HTTP_API_PORT"
	HTTP_HOST          EnvVar = "HTTP_API_HOST"
	HTTP_READ_TIMEOUT  EnvVar = "HTTP_READ_TIMEOUT"
	HTTP_WRITE_TIMEOUT EnvVar = "HTTP_WRITE_TIMEOUT"

	PUBLIC_KEY  EnvVar = "PUBLIC_SIGNING_KEY"
	PRIVATE_KEY EnvVar = "PRIVATE_SIGNING_KEY"

	JWT_SECRET EnvVar = "JWT_SECRET"
)

func GetEnvVar(key EnvVar) (string, error) {
	found, ok := os.LookupEnv(string(key))
	if !ok {
		return "", errors.New("environment variable not found")
	}
	return found, nil
}

func MustGetEnvVar(key EnvVar) string {
	found, err := GetEnvVar(key)
	if err != nil {
		panic(err)
	}
	return found
}
