package secret

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/secret"
	"github.com/coopersmall/subswag/utils"
)

type SecretService struct {
	logger     utils.ILogger
	rsaService rsaService
	secretRepo secretsRepo
	pub        *rsa.PublicKey
	priv       *rsa.PrivateKey
}

func NewSecretService(
	logger utils.ILogger,
	rsaService rsaService,
	secretRepo secretsRepo,
	getSecretsPrivateKey func() (*rsa.PrivateKey, error),
	getSecretsPublicKey func() (*rsa.PublicKey, error),
) *SecretService {
	priv, err := getSecretsPrivateKey()
	if err != nil {
		panic(utils.NewInternalError("failed to get secrets private key", err))
	}
	pub, err := getSecretsPublicKey()
	if err != nil {
		panic(utils.NewInternalError("failed to get secrets public key", err))
	}
	return &SecretService{
		logger:     logger,
		secretRepo: secretRepo,
		rsaService: rsaService,
		pub:        pub,
		priv:       priv,
	}
}

func (s *SecretService) GetSecret(ctx context.Context, secretId secret.SecretID) (secret.Secret, error) {
	secret, err := s.secretRepo.Get(ctx, secretId)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (s *SecretService) CreateSecret(ctx context.Context, value string) (*secret.StoredSecret, error) {
	salt, err := s.generateSalt()
	if err != nil {
		return nil, err
	}
	encrypted, err := s.encrypt(value, salt)
	if err != nil {
		return nil, err
	}
	secret := secret.NewStoredSecret(
		encrypted,
		salt,
	)
	err = s.secretRepo.Create(ctx, secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (s *SecretService) CreateSecretWithId(ctx context.Context, secretId secret.SecretID, value string) (*secret.StoredSecret, error) {
	salt, err := s.generateSalt()
	if err != nil {
		return nil, err
	}
	encrypted, err := s.encrypt(value, salt)
	if err != nil {
		return nil, err
	}
	secret := &secret.StoredSecret{
		SecretBase: secret.SecretBase{
			ID:       secretId,
			Metadata: domain.NewMetadata(),
		},
		Type: secret.SecretTypeStored,
		StoredSecretData: secret.StoredSecretData{
			Value: encrypted,
			Salt:  salt,
		},
	}
	if err != nil {
		return nil, err
	}
	err = s.secretRepo.Create(ctx, secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
func (s *SecretService) UpdateSecret(ctx context.Context, secretId secret.SecretID, value string) (*secret.StoredSecret, error) {
	secret, err := s.secretRepo.Get(ctx, secretId)
	if err != nil {
		return nil, err
	}
	encryptedValue, err := s.encrypt(value, secret.StoredSecretData.Salt)
	if err != nil {
		return nil, err
	}
	secret.StoredSecretData.Value = encryptedValue
	err = s.secretRepo.Update(ctx, secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (s *SecretService) DeleteSecret(ctx context.Context, secretId secret.SecretID) error {
	err := s.secretRepo.Delete(ctx, secretId)
	if err != nil {
		return err
	}
	return nil
}

func (s *SecretService) encrypt(value string, salt []byte) ([]byte, error) {
	salted := value + string(salt)
	return s.rsaService.Encrypt(context.Background(), []byte(salted), s.pub)
}

func (s *SecretService) decrypt(secret *secret.StoredSecret) (string, error) {
	decryptedData, err := s.rsaService.Decrypt(context.Background(), secret.StoredSecretData.Value, s.priv)
	if err != nil {
		return "", err
	}
	salt := string(secret.StoredSecretData.Salt)
	value := string(decryptedData)
	unsalted := value[:len(value)-len(salt)]
	return unsalted, nil
}

const defaultSaltLength = 32

func (s *SecretService) generateSalt() ([]byte, error) {
	salt := make([]byte, defaultSaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	var encoded []byte
	_ = hex.Encode(encoded, salt)
	return encoded, nil
}

type rsaService interface {
	Encrypt(ctx context.Context, data []byte, publicKey *rsa.PublicKey) ([]byte, error)
	Decrypt(ctx context.Context, data []byte, privateKey *rsa.PrivateKey) ([]byte, error)
}
type secretsRepo interface {
	Create(ctx context.Context, secret *secret.StoredSecret) error
	Get(ctx context.Context, secretId secret.SecretID) (*secret.StoredSecret, error)
	Update(ctx context.Context, secret *secret.StoredSecret) error
	Delete(ctx context.Context, secretId secret.SecretID) error
}
