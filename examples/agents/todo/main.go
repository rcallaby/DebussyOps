package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type HandleRequest struct {
	ID      string                 `json:"id"`
	Action  string                 `json:"action"`
	Payload map[string]interface{} `json:"payload"`
}

type HandleResponse struct {
	Status   string                 `json:"status"`
	Response map[string]interface{} `json:"response"`
}

var mu sync.Mutex
var tasks = []map[string]interface{}{}

func main() {
	log.Println("todo agent running :8082")
	http.HandleFunc("/v1/meta", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"name": "todo", "intents": []string{"add_task", "list_tasks"}})
	})

	http.HandleFunc("/v1/handle", func(w http.ResponseWriter, r *http.Request) {
		var hr HandleRequest
		json.NewDecoder(r.Body).Decode(&hr)
		switch hr.Action {
		case "add_task":
			mu.Lock()
			defer mu.Unlock()
			title, _ := hr.Payload["task"].(string)
			id := uuid.New().String()
			t := map[string]interface{}{"id": id, "task": title}
			tasks = append(tasks, t)
			json.NewEncoder(w).Encode(HandleResponse{Status: "ok", Response: map[string]interface{}{"task": t}})
		case "list_tasks":
			json.NewEncoder(w).Encode(HandleResponse{Status: "ok", Response: map[string]interface{}{"tasks": tasks}})
		default:
			w.WriteHeader(400)
			w.Write([]byte("unknown action"))
		}
	})

	log.Fatal(http.ListenAndServe(":8082", nil))
}
