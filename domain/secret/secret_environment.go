package secret

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/utils"
)

type EnvironmentSecretData struct {
	/*
	 @minLength 1
	*/
	Name string `json:"name" validate:"required,min=1"`
}

func (EnvironmentSecretData) IsSecretDataUnion() {}

type EnvironmentSecret struct {
	SecretBase            `json:",inline" tstype:",extends"`
	Type                  SecretType `json:"type" validate:"required,eq=environment" tstype:"'environment'"`
	EnvironmentSecretData `json:",inline" tstype:",extends"`
}

func (e EnvironmentSecret) GetData() SecretDataUnion {
	return e.EnvironmentSecretData
}

func NewEnvironmentSecret(
	name string,
) *EnvironmentSecret {
	return &EnvironmentSecret{
		SecretBase: SecretBase{
			ID:       SecretID(utils.NewID()),
			Metadata: domain.NewMetadata(),
		},
		Type: SecretTypeEnvironment,
		EnvironmentSecretData: EnvironmentSecretData{
			Name: name,
		},
	}
}
