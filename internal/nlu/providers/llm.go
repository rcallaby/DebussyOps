package providers

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/example/orchestrator-template/internal/nlu"
    "github.com/google/uuid"
)

// LLMNLU uses an OpenAI-compatible chat completions API for intelligent intent parsing.
type LLMNLU struct {
    apiKey   string
    model    string
    baseURL  string
    client   *http.Client
    maxRetries int
}

// NewLLMNLU creates a new LLM-based NLU provider.
func NewLLMNLU() *LLMNLU {
    baseURL := getEnv("LLM_BASE_URL", "https://api.openai.com/v1")
    model := getEnv("LLM_MODEL", "gpt-4o-mini") // or gpt-4o, llama-3.1-70b, etc.

    return &LLMNLU{
        apiKey:     os.Getenv("LLM_API_KEY"),
        model:      model,
        baseURL:    strings.TrimSuffix(baseURL, "/"),
        client:     &http.Client{Timeout: 30 * time.Second},
        maxRetries: 2,
    }
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}

// chatMessage represents a message in the chat completion request.
type chatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type chatRequest struct {
    Model          string        `json:"model"`
    Messages       []chatMessage `json:"messages"`
    ResponseFormat map[string]any `json:"response_format,omitempty"`
    Temperature    float64       `json:"temperature"`
}

type chatResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
    Error *struct {
        Message string `json:"message"`
    } `json:"error,omitempty"`
}

// Parse implements the NLU interface using an LLM.
func (l *LLMNLU) Parse(input string) (*nlu.AgentRequest, error) {
    if l.apiKey == "" {
        return nil, errors.New("LLM_API_KEY environment variable is required for LLM NLU")
    }

    systemPrompt := `You are an expert intent classifier for a task automation orchestrator.

Available agents and actions:
- calendar: create_event, list_availability
- todo: add_task, list_tasks

Your job is to convert the user's natural language input into a structured JSON object with exactly these fields:

{
  "agent": "calendar" | "todo",
  "action": string,
  "payload": object
}

Rules:
- Only use the agents and actions listed above.
- Put relevant details into the "payload" object.
- If the input is ambiguous, choose the most reasonable agent/action.
- Always return valid JSON. Do not add explanations or markdown.`

    userPrompt := fmt.Sprintf("User input: %s", input)

    reqBody := chatRequest{
        Model: l.model,
        Messages: []chatMessage{
            {Role: "system", Content: systemPrompt},
            {Role: "user", Content: userPrompt},
        },
        Temperature: 0.1,
    }

    // Use structured outputs if the provider supports it (OpenAI gpt-4o+)
    if strings.Contains(l.model, "gpt-4o") {
        reqBody.ResponseFormat = map[string]any{
            "type": "json_object",
        }
    }

    var lastErr error
    for attempt := 0; attempt <= l.maxRetries; attempt++ {
        result, err := l.callLLM(reqBody)
        if err == nil {
            // Try to parse the response into AgentRequest
            var parsed struct {
                Agent   string                 `json:"agent"`
                Action  string                 `json:"action"`
                Payload map[string]interface{} `json:"payload"`
            }
            if err := json.Unmarshal([]byte(result), &parsed); err == nil {
                return &nlu.AgentRequest{
                    ID:      uuid.New().String(),
                    Agent:   parsed.Agent,
                    Action:  parsed.Action,
                    Payload: parsed.Payload,
                }, nil
            }
            lastErr = fmt.Errorf("failed to parse LLM JSON: %w", err)
        } else {
            lastErr = err
        }

        // Simple repair attempt on retry
        if attempt < l.maxRetries {
            reqBody.Messages = append(reqBody.Messages, chatMessage{
                Role:    "user",
                Content: "The previous response was invalid. Please return ONLY valid JSON matching the schema.",
            })
        }
    }

    return nil, fmt.Errorf("LLM NLU failed after %d attempts: %w", l.maxRetries+1, lastErr)
}

func (l *LLMNLU) callLLM(req chatRequest) (string, error) {
    body, _ := json.Marshal(req)

    httpReq, err := http.NewRequest("POST", l.baseURL+"/v1/chat/completions", bytes.NewReader(body))
    if err != nil {
        return "", err
    }
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+l.apiKey)

    resp, err := l.client.Do(httpReq)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)

    var chatResp chatResponse
    if err := json.Unmarshal(respBody, &chatResp); err != nil {
        return "", err
    }

    if chatResp.Error != nil {
        return "", fmt.Errorf("LLM API error: %s", chatResp.Error.Message)
    }

    if len(chatResp.Choices) == 0 {
        return "", errors.New("no response from LLM")
    }

    return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}