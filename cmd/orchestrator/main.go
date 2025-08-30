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
	"github.com/example/orchestrator-template/internal/router"
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
	log.Infof("Starting Orchestrator (env NLU_PROVIDER=%s)", os.Getenv("NLU_PROVIDER"))

	// choose NLU
	var n nlu.NLU
	switch os.Getenv("NLU_PROVIDER") {
	case "keyword":
		n = providers.NewKeywordNLU()
	default:
		n = providers.NewKeywordNLU()
	}

	// setup orchestrator core
	orch := orchestrator.NewOrchestrator(n, log)

	// register agents (in template, local HTTP agents)
	cal := client.NewAgentClient(cfg.CalendarURL, "")
	todo := client.NewAgentClient(cfg.TodoURL, "")
	torch.RegisterAgent("calendar", cal)
	torch.RegisterAgent("todo", todo)

	// HTTP endpoints
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
		start := time.Now()
		ar, err := torch.ParseInput(q.Input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		res, err := torch.RouteAndExecute(ar)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		_ = log
		_ = start
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	http.ListenAndServe(":8080", nil)
}
