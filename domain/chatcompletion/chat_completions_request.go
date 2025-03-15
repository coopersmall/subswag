package chatcompletion

type ChatCompletionRequest struct {
	Messages     []Message
	Tools        []Tool
	Model        Model
	Temperature  float64
	ResponseType ResponseType
	NextTool     string
}
