# Getting Started with DebussyOps

This guide will help you get DebussyOps running quickly so you can start orchestrating AI agents with natural language.

## What is DebussyOps?

DebussyOps is a **provider-agnostic Master Orchestrator** written in Go. It accepts natural language input, understands the intent using an NLU (Natural Language Understanding) layer, and routes the request to the appropriate specialized agent.

It is designed to be:
- Extensible (easy to add new agents or NLU providers)
- Modular
- Easy to run locally or in Docker

## Prerequisites

- Go 1.20 or higher
- Make (recommended)
- Docker & Docker Compose (optional but recommended)
- An API key from an LLM provider (OpenAI, Groq, Together, etc.) if using the LLM NLU

## Quick Start

### Option 1: Run Locally with Make (Fastest for Development)

```bash
# Clone the repository
git clone https://github.com/rcallaby/DebussyOps.git
cd DebussyOps

# Build the orchestrator and example agents
make build

# Start everything (orchestrator + calendar + todo agents)
make run
```

The orchestrator will be available at **http://localhost:8080**.

### Option 2: Run with Docker Compose (Recommended)

```bash
git clone https://github.com/rcallaby/DebussyOps.git
cd DebussyOps

docker compose up --build
```

This starts three services:
- **Orchestrator** → `http://localhost:8080`
- **Calendar Agent** → `http://localhost:8081`
- **Todo Agent** → `http://localhost:8082`

## Basic Usage

Once the orchestrator is running, send natural language requests:

```bash
curl -X POST http://localhost:8080/v1/query \
  -H "Content-Type: application/json" \
  -d '{"input": "Schedule a meeting with the team tomorrow at 10am"}'
```

**Example responses:**

- Calendar-related request → routes to the `calendar` agent
- Task-related request → routes to the `todo` agent

### Health Check

```bash
curl http://localhost:8080/health
```

## Configuration

DebussyOps is configured using environment variables.

### Basic Configuration

| Variable       | Default                  | Description                     |
|----------------|--------------------------|---------------------------------|
| `CALENDAR_URL` | `http://localhost:8081`  | URL of the calendar agent       |
| `TODO_URL`     | `http://localhost:8082`  | URL of the todo agent           |

### NLU Provider Selection

| Variable        | Default     | Description                                      |
|-----------------|-------------|--------------------------------------------------|
| `NLU_PROVIDER`  | `keyword`   | Set to `llm` to use the intelligent LLM provider |

### LLM Configuration (when using `NLU_PROVIDER=llm`)

| Variable       | Default                        | Description                              |
|----------------|--------------------------------|------------------------------------------|
| `LLM_API_KEY`  | *(required)*                   | Your LLM provider API key                |
| `LLM_MODEL`    | `gpt-4o-mini`                  | Model name                               |
| `LLM_BASE_URL` | `https://api.openai.com/v1`    | Base URL (change for Groq, Ollama, etc.) |

### Example: Using Groq with a fast model

```bash
NLU_PROVIDER=llm \
LLM_API_KEY=gsk_... \
LLM_MODEL=llama-3.1-70b-versatile \
LLM_BASE_URL=https://api.groq.com/openai/v1 \
docker compose up --build
```

## Choosing an NLU Provider

### Keyword NLU (Default)
- Fast and requires no API keys
- Limited to simple keyword matching
- Good for testing and simple demos

### LLM NLU (Recommended)
- Understands complex, natural language
- Much more flexible and intelligent
- Requires an API key

**Recommendation**: Use the **LLM provider** for any real use case.

## Project Structure Overview

```
DebussyOps/
├── cmd/orchestrator/          # Main application entrypoint
├── internal/
│   ├── orchestrator/          # Core routing logic
│   ├── nlu/                   # NLU interface + providers
│   ├── agent/client/          # HTTP client for agents
│   └── plugins/               # Plugin loader (future use)
├── examples/agents/           # Example agent implementations
├── pkg/config/                # Configuration loader
└── docs/                      # Documentation
```

## Common Next Steps

After getting the basic setup running, try these:

1. **Switch to the LLM NLU** (see Configuration section above)
2. **Add a new agent** — See [Adding_Agents.md](Adding_Agents.md)
3. **Create a custom NLU provider** — See [Adding_NLU_Providers.md](Adding_NLU_Providers.md)
4. **Explore the architecture** — See [architecture.md](architecture.md)

## Troubleshooting

### "could not parse intent" error
→ You are using the Keyword NLU. Switch to `NLU_PROVIDER=llm`.

### LLM returns errors or invalid JSON
→ Check that `LLM_API_KEY` is set correctly and the model supports JSON output (most modern models do).

### Port already in use
→ Make sure no other services are running on ports `8080`, `8081`, or `8082`.

### Docker build fails
→ Try rebuilding without cache:
```bash
docker compose build --no-cache
```

## Need Help?

- Read the other documentation files in the `docs/` folder
- Check the main [README.md](../README.md)
- Review the source code — it's intentionally kept clean and well-organized
