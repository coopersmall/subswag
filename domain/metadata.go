package domain

import (
	"encoding/json"
	"time"
)

type Metadata struct {
	CreatedAt time.Time `json:"created_at" validate:"required,nonZeroTime"`
	UpdatedAt time.Time `json:"updated_at,omitempty" validate:"omitempty,nonZeroTime"`
}

func NewMetadata() *Metadata {
	return &Metadata{
		CreatedAt: time.Now().UTC(),
	}
}

func (m *Metadata) Validate() error {
	return Validate(m)
}

func (m *Metadata) UnmarshalJSON(data []byte) error {
	aux := &struct {
		CreatedAt string `json:"created_at" validate:"required"`
		UpdatedAt string `json:"updated_at,omitempty" validate:"omitempty"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if err := Validate(aux); err != nil {
		return err
	}

	var err error
	m.CreatedAt, err = time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return err
	}

	if aux.UpdatedAt != "" {
		m.UpdatedAt, err = time.Parse(time.RFC3339, aux.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Metadata) MarshalJSON() ([]byte, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at,omitempty"`
	}{
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	})
}
