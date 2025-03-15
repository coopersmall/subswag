package ai

const (
	MODEL_LLAMA_3_3_8b = "llama3-8b-8192"
)

func RenderPrompt(prompt string, opts ...func(string) string) string {
	for _, opt := range opts {
		prompt = opt(prompt)
	}
	return prompt
}

func PromptWithJSONSchema(schema string) func(string) string {
	return func(prompt string) string {
		return prompt + `
Please provide a response in the following JSON format:

` + schema + `

Ensure your response is a valid JSON object following this schema.`
	}
}
