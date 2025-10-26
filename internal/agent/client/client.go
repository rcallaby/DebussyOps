package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AgentMessage struct {
	ID      string                 `json:"id"`
	Action  string                 `json:"action"`
	Payload map[string]interface{} `json:"payload"`
}

type AgentResponse struct {
	Status   string                 `json:"status"`
	Response map[string]interface{} `json:"response"`
}

type AgentClient struct {
	BaseURL string
	client  *http.Client
}
// SendMessage sends a message to the agent server
func (a *AgentClient) SendMessage(msg *AgentMessage) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", a.BaseURL+"/v1/message", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
// Handle handles a message from the agent server and returns a response
func (a *AgentClient) Handle(msg *AgentMessage) (*AgentResponse, error) {
	b, _ := json.Marshal(msg)
	req, _ := http.NewRequest("POST", a.BaseURL+"/v1/handle", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var ar AgentResponse
	if err := json.NewDecoder(res.Body).Decode(&ar); err != nil {
		return nil, err
	}
	if ar.Status != "ok" {
		return &ar, fmt.Errorf("agent returned status %s", ar.Status)
	}
	return &ar, nil
}

func NewAgentClient(base string, token string) *AgentClient {
	return &AgentClient{BaseURL: base, client: &http.Client{Timeout: 8 * time.Second}}
}