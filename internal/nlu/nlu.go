package nlu

// AgentRequest created by NLU providers
type AgentRequest struct {
	ID      string                 `json:"id"`
	Agent   string                 `json:"agent"`
	Action  string                 `json:"action"`
	Payload map[string]interface{} `json:"payload"`
}

// NLU interface - plug any provider
type NLU interface {
	Parse(input string) (*AgentRequest, error)
}
