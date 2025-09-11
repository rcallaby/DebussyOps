package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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

var cal = map[string][]map[string]interface{}{}

func main() {
	log.Println("calendar agent running :8081")
	http.HandleFunc("/v1/meta", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"name": "calendar", "intents": []string{"create_event", "list_availability"}})
	})

	http.HandleFunc("/v1/handle", func(w http.ResponseWriter, r *http.Request) {
		var hr HandleRequest
		json.NewDecoder(r.Body).Decode(&hr)
		switch hr.Action {
		case "create_event":
			title, _ := hr.Payload["title"].(string)
			id := uuid.New().String()
			e := map[string]interface{}{"id": id, "title": title, "time": time.Now().Add(24 * time.Hour).Format(time.RFC3339)}
			cal["default"] = append(cal["default"], e)
			json.NewEncoder(w).Encode(HandleResponse{Status: "ok", Response: map[string]interface{}{"event": e}})
		case "list_availability":
			json.NewEncoder(w).Encode(HandleResponse{Status: "ok", Response: map[string]interface{}{"avail": []string{time.Now().Add(48 * time.Hour).Format(time.RFC3339)}})})
		default:
			w.WriteHeader(400)
			w.Write([]byte("unknown action"))
		}
	))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
