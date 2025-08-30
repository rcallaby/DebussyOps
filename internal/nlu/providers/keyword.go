package providers

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/example/orchestrator-template/internal/nlu"
)

type KeywordNLU struct{}

func NewKeywordNLU() *KeywordNLU { return &KeywordNLU{} }

func containsAny(input string, list []string) bool {
	in := strings.ToLower(input)
	for _, w := range list {
		if strings.Contains(in, w) { return true }
	}
	return false
}

func (k *KeywordNLU) Parse(input string) (*nlu.AgentRequest, error) {
	if containsAny(input, []string{"meeting", "schedule", "calendar"}) {
		return &nlu.AgentRequest{ID: uuid.New().String(), Agent: "calendar", Action: "create_event", Payload: map[string]interface{}{"title": input}}, nil
	}
	if containsAny(input, []string{"task", "todo", "remind"}) {
		return &nlu.AgentRequest{ID: uuid.New().String(), Agent: "todo", Action: "add_task", Payload: map[string]interface{}{"task": input}}, nil
	}
	return nil, errors.New("could not parse intent")
}
