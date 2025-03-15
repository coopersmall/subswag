package chatcompletion

import "encoding/json"

type ChatCompletionResponse struct {
	Content   json.RawMessage
	ToolCalls []ToolCall
}
