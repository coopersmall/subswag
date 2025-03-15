package chatcompletion

import "errors"

type ResponseType string

const (
	ResponseTypeText       ResponseType = "text"
	ResponseTypeJSONObject ResponseType = "json_object"
)

func ResponseTypes() []ResponseType {
	return []ResponseType{
		ResponseTypeText,
		ResponseTypeJSONObject,
	}
}

func (r ResponseType) String() string {
	return string(r)
}

func (r ResponseType) Validate() error {
	responseTypes := ResponseTypes()
	for _, responseType := range responseTypes {
		if responseType == r {
			return nil
		}
	}
	return errors.New("invalid response type")
}
