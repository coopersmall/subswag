package main

import (
	"os"

	"github.com/coopersmall/subswag/env"
	"github.com/coopersmall/subswag/utils"
)

func main() {
	if _, err := os.Stat(utils.ENV_FILE); !os.IsNotExist(err) {
		os.Remove(utils.ENV_FILE)
	}

	envFile, err := os.Create(utils.ENV_FILE)
	if err != nil {
		panic(err)
	}

	_, err = envFile.WriteString(localEnv())
	if err != nil {
		panic(err)
	}
}

func localEnv() string {
	requestPublicKey, requestPrivateKey, err := utils.GenerateRSAKeyPair()
	if err != nil {
		panic(err)
	}

	secretsPublicKey, secretsPrivateKey, err := utils.GenerateRSAKeyPair()
	if err != nil {
		panic(err)
	}

	jwtSecret := "some-secure-secret"

	// We need to generate an actual secret here
	return `` + string(env.API_URL) + `="localhost:9000"
` + string(env.API_TIMEOUT) + `="10s"

` + string(env.POSTGRES_URL) + `="postgres://user:password@localhost:5432/subswag?sslmode=disable"
` + string(env.OPENAI_API_KEY) + `="some-openai-key"

JWT_SIGNING_ALGO="HS256"
JWT_DEFAULT_EXP="24h"
` + string(env.JWT_SIGNING_KEY) + `="` + jwtSecret + `"

` + string(env.REQUESTS_PUBLIC_KEY) + `="` + utils.EncodePublicKeyToPEM(requestPublicKey) + `"

` + string(env.REQUESTS_PRIVATE_KEY) + `="` + utils.EncodePrivateKeyToPEM(requestPrivateKey) + `"

` + string(env.SECRETS_PUBLIC_KEY) + `="` + utils.EncodePublicKeyToPEM(secretsPublicKey) + `"

` + string(env.SECRETS_PRIVATE_KEY) + `="` + utils.EncodePrivateKeyToPEM(secretsPrivateKey) + `"
`
}
