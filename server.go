package main

import (
    "encoding/json"
    "net/http"
    "os"
    "time"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{
        "message":   "pong",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    http.HandleFunc("/ping", pingHandler)

    port := os.Getenv("PORT")
    if port == "" {
        port = "10000" // fallback if not set by Render
    }

    http.ListenAndServe(":"+port, nil)
}