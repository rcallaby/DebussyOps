# DebussyOps


A provider-agnostic, extensible Master Orchestrator template in Go that coordinates modular sub-agents with plugin support and natural language understanding for both personal and business automation.

This repository is a scaffold for an extensible AI Orchestrator. Use it as a starting point to build domain-specific orchestrators (personal life, business automation, etc.).

Features:
- Modular agent interface
- Provider-agnostic NLU abstraction
- Plugin loader stub (for future dynamic modules)
- Config via Viper
- Structured logging (zap)
- Example agents: calendar & todo

## Quick start (local)

Requirements: Go 1.20+, Make

```bash
make build
make run
```

This starts the orchestrator on :8080 and two example agents on :8081 and :8082.

API:
- POST /v1/query { "input": "..." }
- GET  /health


