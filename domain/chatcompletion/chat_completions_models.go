package chatcompletion

import "errors"

type Model string

const (
	ModelGPT4o Model = "4o"
)

func Models() []Model {
	return []Model{
		ModelGPT4o,
	}
}

func (m Model) String() string {
	return string(m)
}

func (m Model) Validate() error {
	models := Models()
	for _, model := range models {
		if model == m {
			return nil
		}
	}
	return errors.New("invalid model")
}
