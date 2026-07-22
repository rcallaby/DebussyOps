package main

import (
    "encoding/json"
    "net/http"
    "os"
    "time"

    "go.uber.org/zap"

    "github.com/example/orchestrator-template/internal/nlu"
    "github.com/example/orchestrator-template/internal/nlu/providers"
    "github.com/example/orchestrator-template/internal/orchestrator"
    "github.com/example/orchestrator-template/internal/agent/client"
    "github.com/example/orchestrator-template/pkg/config"
)

type QueryRequest struct {
    Input string `json:"input"`
}

func main() {
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    log := logger.Sugar()

    cfg := config.Load()
    log.Infof("Starting DebussyOps Orchestrator (NLU_PROVIDER=%s)", os.Getenv("NLU_PROVIDER"))

    // === NLU Selection ===
    var nluProvider nlu.NLU
    switch os.Getenv("NLU_PROVIDER") {
    case "llm", "openai", "gpt":
        nluProvider = providers.NewLLMNLU()
        log.Info("Using LLM-powered NLU")
    default:
        nluProvider = providers.NewKeywordNLU()
        log.Info("Using Keyword NLU (fallback)")
    }

    // Core orchestrator
    orch := orchestrator.NewOrchestrator(nluProvider, log)

    // Register agents
    cal := client.NewAgentClient(cfg.CalendarURL, "")
    todo := client.NewAgentClient(cfg.TodoURL, "")
    orch.RegisterAgent("calendar", cal)
    orch.RegisterAgent("todo", todo)

    // HTTP server
    http.HandleFunc("/v1/query", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        var q QueryRequest
        if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("invalid JSON"))
            return
        }

        ar, err := orch.ParseInput(q.Input)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }

        res, err := orch.RouteAndExecute(ar)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(res)
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    log.Info("Orchestrator listening on :8080")
    http.ListenAndServe(":8080", nil)
}
