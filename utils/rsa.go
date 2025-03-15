package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func MustGetRSAPublicKey() *rsa.PublicKey {
	key, err := GetRSAPublicKey()
	if err != nil {
		panic(err)
	}
	return key
}

func GetRSAPublicKey() (*rsa.PublicKey, error) {
	found, err := GetEnvVar(PUBLIC_KEY)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(found))
	if block == nil {
		return nil, err
	}

	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		return key.(*rsa.PublicKey), nil
	}
	return nil, err
}

func MustGetRSASigningKey() *rsa.PrivateKey {
	key, err := GetRSASigningKey()
	if err != nil {
		panic(err)
	}
	return key
}

func GetRSASigningKey() (*rsa.PrivateKey, error) {
	found, err := GetEnvVar(PRIVATE_KEY)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(found))
	if block == nil {
		return nil, err
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

func GenerateRSAKeyPair() (*rsa.PublicKey, *rsa.PrivateKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return &privKey.PublicKey, privKey, nil
}

func EncodePublicKeyToPEM(key *rsa.PublicKey) string {
	pubASN1, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		panic(err)
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	}
	return string(pem.EncodeToMemory(block))
}

func EncodePrivateKeyToPEM(key *rsa.PrivateKey) string {
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return string(pem.EncodeToMemory(block))
}
