package encryption

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/coopersmall/subswag/utils"
)

type RSAService struct {
	logger utils.ILogger
}

func NewRSAService(logger utils.ILogger) *RSAService {
	return &RSAService{
		logger: logger,
	}
}

func (r *RSAService) Encrypt(
	ctx context.Context,
	data []byte,
	publicKey *rsa.PublicKey,
) ([]byte, error) {
	return rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		data,
		nil,
	)
}

func (r *RSAService) Decrypt(
	ctx context.Context,
	data []byte,
	privateKey *rsa.PrivateKey,
) ([]byte, error) {
	return rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privateKey,
		data,
		nil,
	)
}
