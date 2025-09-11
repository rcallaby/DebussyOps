# Architecture

> Template: **orchestrator-template** (Go)
> Purpose: a provider-agnostic, extensible orchestrator that routes natural-language requests to modular sub-agents over HTTP.

This document describes the runtime flow, components, data contracts, configuration, and known limitations of the scaffold. It reflects the actual files and code currently in the template—no features are described that aren’t present.

---

## High-level Overview

* Clients send a POST request to the **Orchestrator** (`/v1/query`) with free-text input.
* The orchestrator calls an **NLU provider** to parse the text into a structured `AgentRequest`.
* The orchestrator looks up a registered **Agent Client** by name (e.g., `calendar`, `todo`).
* It forwards a JSON payload to the agent’s HTTP endpoint (`/v1/handle`).
* The orchestrator returns the agent’s `response` object **as-is** to the caller.

```
Client
  │  POST /v1/query {"input": "..."}
  ▼
Orchestrator (cmd/orchestrator + internal/orchestrator)
  │  Parse(input) -> AgentRequest {Agent, Action, Payload}
  ▼
NLU (internal/nlu + providers/keyword)
  │
  ▼
Agent Client (internal/agent/client)  ──►  Agent HTTP Service (/v1/handle)
  │                                           ▲
  ▼                                           │
returns {status:"ok", response:{...}}  ◄──────┘
  │
  ▼
Orchestrator returns only the inner `response` JSON map to caller
```

---

## Repository Layout 

```
.
├─ cmd/
│  └─ orchestrator/
│     └─ main.go
├─ internal/
│  ├─ orchestrator/
│  │  └─ orchestrator.go
│  ├─ nlu/
│  │  ├─ nlu.go
│  │  └─ providers/
│  │     └─ keyword.go
│  ├─ agent/
│  │  ├─ agent.go                 # unused in current flow; for plugin-style internal agents
│  │  └─ client/
│  │     └─ client.go             # HTTP client to external agent services
│  └─ plugins/
│     └─ plugin_loader.go         # stub (comments only)
├─ pkg/
│  └─ config/
│     └─ config.go                # minimal env loader
├─ examples/
│  └─ agents/
│     ├─ calendar/
│     │  └─ main.go               # HTTP agent :8081
│     └─ todo/
│        └─ main.go               # HTTP agent :8082
├─ README.md
├─ Makefile
├─ docker-compose.yml
├─ Dockerfile                     # builds only orchestrator
└─ go.mod
```

---

## Components & Responsibilities

### 3.1 Orchestrator entrypoint

**File:** `cmd/orchestrator/main.go`

* Wires together:

  * logger (`zap`)
  * config (env) via `pkg/config/config.go`
  * NLU provider selection via `NLU_PROVIDER` env (only `keyword` and default→`keyword` supported)
  * agent registry using HTTP base URLs from config
* Exposes:

  * `POST /v1/query` – main orchestration endpoint
  * `GET /health` – basic liveness check
* Flow in `/v1/query`:

  1. Decode `{"input": "<free text>"}`.
  2. `ParseInput` via NLU → `AgentRequest`.
  3. `RouteAndExecute` → call agent `/v1/handle`.
  4. Return **only** the agent’s `response` map as the response body.
* Implementation detail:

  * It sets local variables `start` and `log` and assigns to `_` to avoid unused warnings; there is no current latency metric.

---

### 3.2 Orchestrator core

**File:** `internal/orchestrator/orchestrator.go`

* Holds:

  * `nlu` (implements `nlu.NLU`)
  * `log` (`*zap.SugaredLogger`)
  * `agents` map: `name → *agentclient.AgentClient`
* Methods:

  * `RegisterAgent(name string, c *agentclient.AgentClient)`
  * `ParseInput(input string) (*nlu.AgentRequest, error)` – delegates to NLU
  * `RouteAndExecute(ar *nlu.AgentRequest) (map[string]interface{}, error)` – constructs agent message and forwards; returns `res.Response`
* No persistence, retries, or circuit breaking implemented.

---

### 3.3 NLU abstraction & default provider

**Files:**

* `internal/nlu/nlu.go`

* `internal/nlu/providers/keyword.go`

* `nlu.NLU` interface:

  ```go
  type NLU interface {
      Parse(input string) (*AgentRequest, error)
  }
  ```

* `AgentRequest` produced by NLU:

  ```go
  type AgentRequest struct {
      ID      string
      Agent   string
      Action  string
      Payload map[string]interface{}
  }
  ```

* **KeywordNLU** (baseline):

  * If input contains `meeting`, `schedule`, or `calendar` → `Agent:"calendar"`, `Action:"create_event"`.
  * If input contains `task`, `todo`, or `remind` → `Agent:"todo"`, `Action:"add_task"`.
  * Otherwise returns an error (“could not parse intent”).

---

### 3.4 Agent client (HTTP)

**File:** `internal/agent/client/client.go`

* Contract used to call external agent services:

  * Request to `/v1/handle`:

    ```go
    type AgentMessage struct {
      ID      string                 `json:"id"`
      Action  string                 `json:"action"`
      Payload map[string]interface{} `json:"payload"`
    }
    ```
  * Expected response:

    ```go
    type AgentResponse struct {
      Status   string                 `json:"status"`   // expects "ok" on success
      Response map[string]interface{} `json:"response"` // arbitrary payload
    }
    ```
* HTTP behavior:

  * POSTs JSON with an 8s client timeout.
  * Treats any `Status` value other than `"ok"` as an error.

---

### 3.5 Example agents (HTTP microservices)

#### Calendar Agent

**File:** `examples/agents/calendar/main.go`
**Port:** `:8081`
**Endpoints:**

* `GET /v1/meta` → `{"name":"calendar","intents":["create_event","list_availability"]}`
* `POST /v1/handle`:

  * `action=="create_event"` → creates a fake event for “tomorrow” and returns it.
  * `action=="list_availability"` → returns a single future timestamp.

> **Accuracy note:** The `list_availability` handler has an **extra closing parenthesis**) in the `json.NewEncoder(...).Encode(...)` call. It must be fixed for the file to compile.

#### Todo Agent

**File:** `examples/agents/todo/main.go`
**Port:** `:8082`
**Endpoints:**

* `GET /v1/meta` → `{"name":"todo","intents":["add_task","list_tasks"]}`
* `POST /v1/handle`:

  * `action=="add_task"` → appends to in-memory slice and returns the task.
  * `action=="list_tasks"` → returns all tasks.

---

### 3.6 Plugin loader (stub)

**File:** `internal/plugins/plugin_loader.go`
Contains comments outlining possible approaches (Go plugins, gRPC/HTTP) but **no executable code**.

---

### 3.7 Configuration

**File:** `pkg/config/config.go`
Minimal env loader (no Viper in this scaffold):

* `CALENDAR_URL` (default `http://localhost:8081`)
* `TODO_URL` (default `http://localhost:8082`)

**Runtime selection:**

* `NLU_PROVIDER` env var read in `cmd/orchestrator/main.go`

  * Supported values: `"keyword"` (any other value falls back to `"keyword"` in current code).

---

## Public HTTP API (current)

### Orchestrator

* `POST /v1/query`

  * **Request**:

    ```json
    { "input": "Schedule a meeting tomorrow at 10am" }
    ```
  * **Response**: the **agent’s `response` object only**, e.g.:

    ```json
    { "event": { "id": "...", "title": "Schedule a meeting tomorrow at 10am", "time": "2025-08-30T..." } }
    ```
  * **Errors**:

    * 400 on invalid JSON or NLU parse error (plain text body).
    * 500 on downstream agent errors (plain text body).

* `GET /health` → `"ok"` (text).

### Agents (expected contract)

* `GET /v1/meta` → informational (not used by orchestrator).
* `POST /v1/handle`

  * **Request** (`AgentMessage`):

    ```json
    { "id": "uuid", "action": "create_event", "payload": { "title": "..." } }
    ```
  * **Response** (`AgentResponse`):

    ```json
    { "status": "ok", "response": { /* agent-defined */ } }
    ```

---

## Build & Runtime

* **Makefile**

  * `make build` → builds three binaries into `./bin`: `orchestrator`, `calendar`, `todo`
  * `make run` → launches all three locally with small delays
* **Dockerfile (root)** builds only the orchestrator binary.
* **docker-compose.yml**

  * Declares three services (`orchestrator`, `calendar`, `todo`) and maps ports 8080/8081/8082.
  * Sets `NLU_PROVIDER=keyword` for orchestrator.

---

## Control Flow (Step-by-step)

1. Client → Orchestrator `/v1/query` (POST).
2. Orchestrator → NLU (`Parse(input)`) → `AgentRequest{Agent, Action, Payload}`.
3. Orchestrator selects registered agent client by `AgentRequest.Agent`.
4. Orchestrator → Agent `/v1/handle` with `AgentMessage{ID, Action, Payload}`.
5. Agent returns `AgentResponse{status:"ok", response:{...}}`.
6. Orchestrator responds to client with the inner `response` map (not wrapped).

---

## Extensibility Points (present, not hypothetical)

* **Add NLU providers**: implement `nlu.NLU` and update the `switch` in `cmd/orchestrator/main.go` to instantiate it (and optionally extend env selection).
* **Add agents (HTTP services)**:

  * Implement `/v1/handle` responding with the `AgentResponse` contract.
  * Instantiate a new `AgentClient` in `cmd/orchestrator/main.go`.
  * Register with `RegisterAgent("name", client)`.
* **(Optional) Internal agents**:

  * `internal/agent/agent.go` defines an interface that could be used if you later embed agents in-process (not used by current orchestrator path).

---

## 8) Dependencies

* Go `1.20`
* `go.uber.org/zap` for structured logging
* `github.com/google/uuid` for IDs

No database, message queue, or metrics libraries are included in the scaffold.

---

## Non-Goals / Not Implemented (by design in this scaffold)

* Authentication/Authorization
* Persistent storage
* Auto-discovery of agents or meta negotiation
* Retries, backoff, or circuit breakers
* Metrics/Tracing
* Rich error envelopes (errors are plain text)
* Production config management (the loader is intentionally minimal)
* Agent Dockerfiles (see note under §5)

---

## 10) Known Gaps to Fix (for clean build & compose)

1. **Unused / missing import**

   * `cmd/orchestrator/main.go` imports `github.com/example/orchestrator-template/internal/router`, but no such package exists in the scaffold and it’s not used.
   * **Fix:** remove that import.

2. **Calendar agent syntax error**

   * In `examples/agents/calendar/main.go`, the `list_availability` response has an extra `)` causing a compile error.
   * **Fix:** remove the superfluous parenthesis in the `json.NewEncoder(...).Encode(...)` call.

3. **Dockerfiles for agents**

   * `docker-compose.yml` expects Dockerfiles in `./examples/agents/calendar` and `./examples/agents/todo`.
   * **Fix:** add minimal Dockerfiles in those directories or adjust compose to build/run them differently.

---

## Security & Production Hardening (not present; recommended later)

* TLS, auth (mTLS/JWT), input validation beyond keyword checks
* Request timeouts and context propagation end-to-end
* Observability: Prometheus metrics, OpenTelemetry tracing, structured error contracts
* Persistence (SQLite/Postgres) and idempotency keys for agent actions
* Config via Viper, environment overlays, and secrets management
* Resilience patterns: retries, exponential backoff, circuit breakers
* Agent versioning and compatibility checks via `/v1/meta`

---

## Example Requests

Create an event (keyword NLU → calendar):

```bash
curl -s -X POST http://localhost:8080/v1/query \
  -H "Content-Type: application/json" \
  -d '{"input":"schedule a meeting tomorrow"}'
```

Add a task (keyword NLU → todo):

```bash
curl -s -X POST http://localhost:8080/v1/query \
  -H "Content-Type: application/json" \
  -d '{"input":"add task: buy milk"}'
```

