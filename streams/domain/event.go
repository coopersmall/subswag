package domain

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/utils"
)

type EventID utils.ID

type Event[Data any] struct {
	ID       EventID          `json:"id"`
	Stream   string           `json:"stream"`
	Data     Data             `json:"data"`
	Metadata *domain.Metadata `json:"metadata"`
}

func NewEvent[Data any](
	stream string,
	data Data,
) *Event[Data] {
	return &Event[Data]{
		ID:       EventID(utils.NewID()),
		Stream:   stream,
		Data:     data,
		Metadata: domain.NewMetadata(),
	}
}
