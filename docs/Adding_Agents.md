# Adding New Agents to the Orchestrator

This guide explains how to add additional agents to the orchestrator so the system can handle more domains (finance, travel, healthcare, personal tasks, etc.). The design is intentionally modular so new agents can be dropped in with minimal changes.

---

## 1. Project Structure (Relevant Parts)

```
internal/
  agents/
    base.go          # Agent interface definition
    weather.go       # Example agent
    calendar.go      # Example agent
  nlu/
    nlu.go           # Input parser → AgentRequest
  orchestrator/
    orchestrator.go  # Routes requests to agents
```

---

## 2. The Agent Interface

All agents implement a common interface defined in `internal/agents/base.go`:

```go
package agents

type Agent interface {
    // Name should return a unique identifier for the agent
    Name() string

    // CanHandle determines if this agent is able to process the request
    CanHandle(req *AgentRequest) bool

    // Handle performs the agent's action and returns the response
    Handle(req *AgentRequest) (string, error)
}
```

* **`Name()`** → Used by the orchestrator for routing/debugging.
* **`CanHandle()`** → Filters which agent can respond to a given request.
* **`Handle()`** → Executes the logic (API call, computation, etc.).

---

## 3. Creating a New Agent

1. **Create a new file** under `internal/agents/` (e.g., `todo.go`).
2. **Implement the `Agent` interface**:

```go
package agents

import "fmt"

type TodoAgent struct{}

func (a *TodoAgent) Name() string {
    return "todo"
}

func (a *TodoAgent) CanHandle(req *AgentRequest) bool {
    return req.Intent == "todo" // comes from NLU
}

func (a *TodoAgent) Handle(req *AgentRequest) (string, error) {
    // Example: Add a task
    task := req.Parameters["task"]
    return fmt.Sprintf("Task '%s' added to your list ✅", task), nil
}
```

3. **Register your agent** in the orchestrator:

```go
// internal/orchestrator/orchestrator.go

import "myproject/internal/agents"

func (o *Orchestrator) RegisterDefaultAgents() {
    o.Register(&agents.WeatherAgent{})
    o.Register(&agents.CalendarAgent{})
    o.Register(&agents.TodoAgent{})   // <-- add here
}
```

---

## 4. Extending the NLU

The **NLU** maps natural language → `AgentRequest`.

Example `AgentRequest`:

```go
type AgentRequest struct {
    Intent     string
    Parameters map[string]string
}
```

* If your agent expects a new intent (`"todo"` in our example), extend the NLU (`internal/nlu/nlu.go`) to parse it.

Example (simplified keyword parser):

```go
if strings.Contains(input, "add task") {
    return &AgentRequest{
        Intent: "todo",
        Parameters: map[string]string{
            "task": extractTask(input),
        },
    }, nil
}
```

If using an **LLM-based NLU**, just add `"todo"` to the prompt specification so it can return this intent.

---

## 5. Testing Your Agent

* Write a unit test under `internal/agents/`:

```go
func TestTodoAgent(t *testing.T) {
    agent := &TodoAgent{}
    req := &AgentRequest{
        Intent: "todo",
        Parameters: map[string]string{"task": "Buy milk"},
    }

    resp, err := agent.Handle(req)
    if err != nil {
        t.Fatal(err)
    }

    expected := "Task 'Buy milk' added to your list ✅"
    if resp != expected {
        t.Errorf("expected %q, got %q", expected, resp)
    }
}
```

Run tests:

```sh
go test ./internal/agents/...
```

---

## 6. Best Practices

* Keep **agents focused** on a single domain (e.g., TravelAgent, FinanceAgent).
* Avoid mixing unrelated responsibilities in one agent.
* Validate input inside `Handle()` — never assume parameters are present.
* Document your agent’s purpose at the top of its `.go` file.
* If your agent relies on **external APIs**, place credentials/configs in `config/` or environment variables (never hardcode keys).

---

## 7. Summary

To add a new agent:

1. Implement `Agent` in a new file.
2. Register it in the orchestrator.
3. Extend the NLU to recognize the new intent.
4. Write tests to validate functionality.


