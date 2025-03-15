package secret

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/utils"
)

type StoredSecretData struct {
	Salt  []byte `json:"salt" validate:"required,min=1"`
	Value []byte `json:"value" validate:"required,min=1"`
}

func (StoredSecretData) IsSecretDataUnion() {}

type StoredSecret struct {
	SecretBase       `json:",inline" tstype:",extends"`
	Type             SecretType `json:"type" validate:"required, eq=stored" tstype:"'stored'"`
	StoredSecretData `json:",inline" tstype:",extends"`
}

func (s StoredSecret) GetData() SecretDataUnion {
	return s.StoredSecretData
}

func NewStoredSecret(
	value []byte,
	salt []byte,
) *StoredSecret {
	return &StoredSecret{
		SecretBase: SecretBase{
			ID:       SecretID(utils.NewID()),
			Metadata: domain.NewMetadata(),
		},
		Type: SecretTypeStored,
		StoredSecretData: StoredSecretData{
			Salt:  salt,
			Value: value,
		},
	}
}
