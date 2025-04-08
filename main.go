package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Stat struct {
	Website string         `json:"website"`
	Data    map[string]any `json:"data"`
}

type Server struct {
	AllowedWebsites map[string]bool
	Stats           map[string][]Stat
	mu              sync.Mutex
}

func NewServer() *Server {
	return &Server{
		AllowedWebsites: map[string]bool{
			"example.com": true,
			"mysite.org":  true,
			"another.io":  true,
		},
		Stats: make(map[string][]Stat),
	}
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var stat Stat
	if err := json.NewDecoder(r.Body).Decode(&stat); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(stat)

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.AllowedWebsites[stat.Website] {
		http.Error(w, "Website not allowed", http.StatusForbidden)
		return
	}

	s.Stats[stat.Website] = append(s.Stats[stat.Website], stat)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"success"}`))
}

func main() {
	server := NewServer()

	http.HandleFunc("/stats", server.handleStats)

	log.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
