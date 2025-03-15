//ts:ignore
package secret

import (
	"github.com/coopersmall/subswag/domain"
)

type SecretType string

const (
	SecretTypeEnvironment SecretType = "environment"
	SecretTypeStored      SecretType = "stored"
)

type SecretDataUnion interface {
	IsSecretDataUnion()
}

type Secret interface {
	GetID() SecretID
	GetData() SecretDataUnion
	GetMetadata() *domain.Metadata
}
