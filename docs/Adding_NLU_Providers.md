# Guide: Adding a New NLU Provider

The Master Orchestrator (`lifeos-orchestrator`) supports a **pluggable Natural Language Understanding (NLU) system**.
This allows contributors to integrate different LLMs or NLP services (OpenAI, Anthropic, Cohere, Rasa, HuggingFace, local LLMs, etc.) without modifying core orchestrator logic.

---

## 1. Understand the NLU Interface

All NLU providers must implement the `NLU` interface defined in `internal/nlu/nlu.go`:

```go
// NLU defines the contract for any Natural Language Understanding provider
type NLU interface {
    Parse(input string) (*AgentRequest, error)
}
```

* **`Parse(input string)`**: Takes a raw user string (e.g., `"Remind me to pay my bills tomorrow"`) and returns a structured `AgentRequest` containing:

  * `Agent` → target agent (e.g., `calendar`, `finance`)
  * `Action` → what to do (e.g., `create_event`, `add_expense`)
  * `Params` → extra structured details

---

## 2. Create Your Provider

Create a new file inside `internal/nlu/providers/`.
For example, if you are adding an **OpenAI provider**:

```
internal/nlu/providers/openai.go
```

Example skeleton:

```go
package providers

import (
    "context"
    "fmt"
    "lifeos-orchestrator/internal/nlu"
)

type OpenAINLU struct {
    APIKey string
}

func NewOpenAINLU(apiKey string) *OpenAINLU {
    return &OpenAINLU{APIKey: apiKey}
}

// Parse implements the NLU interface
func (o *OpenAINLU) Parse(input string) (*nlu.AgentRequest, error) {
    // Call your LLM/NLP provider here
    // Convert free text → AgentRequest
    // Example: "Schedule meeting tomorrow" → {Agent:"calendar", Action:"create_event", Params:{...}}

    // Placeholder return until integration
    return &nlu.AgentRequest{
        Agent:  "calendar",
        Action: "create_event",
        Params: map[string]string{"title": "Meeting", "time": "tomorrow"},
    }, nil
}
```

---

## 3. Register Provider in Orchestrator

In `cmd/orchestrator/main.go`, configure which NLU provider to load:

```go
var nluProvider nlu.NLU

switch os.Getenv("NLU_PROVIDER") {
case "openai":
    nluProvider = providers.NewOpenAINLU(os.Getenv("OPENAI_API_KEY"))
case "cohere":
    nluProvider = providers.NewCohereNLU(os.Getenv("COHERE_API_KEY"))
default:
    nluProvider = &nlu.DummyNLU{}
}
```

This allows swapping providers by changing an environment variable (`NLU_PROVIDER=openai`).

---

## 4. Testing Your Provider

Run the orchestrator with your provider:

```bash
NLU_PROVIDER=openai OPENAI_API_KEY=sk-123 make run
```

Send a test request:

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"input":"Schedule a team meeting tomorrow at 3 PM"}'
```

Expected output:

```json
{
  "agent": "calendar",
  "action": "create_event",
  "params": {
    "title": "team meeting",
    "time": "tomorrow 3 PM"
  }
}
```

---

## 5. Contributing New Providers

When submitting a new provider:

* Place implementation in `internal/nlu/providers/`.
* Keep API keys/config in environment variables (no hardcoding).
* Include example usage in the PR.
* Add tests in `internal/nlu/providers/<provider>_test.go`.


