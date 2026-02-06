package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Request represents a standard GraphQL request body
type Request struct {
	Query string `json:"query"`
}

// Response represents a standard GraphQL response
type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []Error     `json:"errors,omitempty"`
}

type Error struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/graphql", handleGraphQL)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := ":8080"
	fmt.Printf("Backend server running on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleGraphQL(w http.ResponseWriter, r *http.Request) {
	// CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Minimal "GraphQL" Parser & Resolver
	// Supports: query { hello, time, user { name } }
	query := strings.TrimSpace(req.Query)
	query = strings.ReplaceAll(query, "\n", " ")
	
	response := Response{Data: make(map[string]interface{})}
	data := response.Data.(map[string]interface{})

	// Very naive parsing (scanning for keywords)
	if strings.Contains(query, "hello") {
		data["hello"] = "Hello from Go Backend!"
	}
	if strings.Contains(query, "version") {
		data["version"] = "1.0.0"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
