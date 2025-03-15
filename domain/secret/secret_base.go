package secret

import "github.com/coopersmall/subswag/domain"

type SecretBase struct {
	ID SecretID `json:"id" validate:"required,gt=0"`
	//tygo:emit type: SecretType
	Metadata *domain.Metadata `json:"metadata" validate:"required" tstype:"Metadata"`
}

func (s *SecretBase) GetID() SecretID {
	return s.ID
}

func (s *SecretBase) GetMetadata() *domain.Metadata {
	return s.Metadata
}
