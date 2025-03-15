package integrations

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/utils"
)

type IntegrationID utils.ID

type Integration struct {
	ID          IntegrationID       `json:"id" validate:"required,gt=0"`
	Name        string              `json:"name" validate:"required" tstype:"string"`
	Description string              `json:"description" validate:"required" tstype:"string"`
	Actions     []IntegrationAction `json:"actions" validate:"required" tstype:"Array<IntegrationAction>"`
	Metadata    *domain.Metadata    `json:"metadata" validate:"required" tstype:"Metadata"`
}

type IntegrationActionType string

const (
	IntegrationActionTypeJS IntegrationActionType = "js"
)

type IntegrationAction interface {
	GetBase() IntegrationActionDataBase
	GetData() IntegrationActionDataUnion
}

func NewIntegrationAction(
	data IntegrationActionDataUnion,
) IntegrationAction {
	switch d := data.(type) {
	case IntegrationJSAction:
		return NewIntegrationJSAction(
			d.InputSchema,
			d.Script,
			d.OutputSchema,
		)
	default:
		return nil
	}
}

type IntegrationActionDataUnion interface {
	IsIntegrationData()
}

type IntegrationActionDataBase struct {
	ID   IntegrationID         `json:"id" validate:"required,gt=0"`
	Type IntegrationActionType `json:"type" validate:"required" tstype:"IntegrationType"`
}

func (i IntegrationActionDataBase) GetBase() IntegrationActionDataBase {
	return i
}

func newIntegrationBaseData(t IntegrationActionType) IntegrationActionDataBase {
	return IntegrationActionDataBase{
		ID:   IntegrationID(utils.NewID()),
		Type: t,
	}
}

type IntegrationJSAction struct {
	IntegrationActionDataBase `json:",inline" validate:"required" tstype:",extends"`
	IntegrationJSActionData   `json:",inline" validate:"required" tstype:",extends"`
}

func (i IntegrationJSAction) GetData() IntegrationActionDataUnion {
	return i.IntegrationJSActionData
}

type IntegrationJSActionData struct {
	InputSchema  string `json:"inputSchema" validate:"required" tstype:"string"`
	Script       string `json:"script" validate:"required" tstype:"string"`
	OutputSchema string `json:"outputSchema" validate:"required" tstype:"string"`
}

func (IntegrationJSActionData) IsIntegrationData() {}

func NewIntegrationJSAction(
	inputSchema string,
	script string,
	outputSchema string,
) *IntegrationJSAction {
	return &IntegrationJSAction{
		IntegrationActionDataBase: newIntegrationBaseData(IntegrationActionTypeJS),
		IntegrationJSActionData: IntegrationJSActionData{
			InputSchema:  inputSchema,
			Script:       script,
			OutputSchema: outputSchema,
		},
	}
}
