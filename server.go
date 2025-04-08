package main

import (
    "encoding/json"
    "io"
    "net/http"
    "os"
    "strconv"
    "time"
)

const (
    maxUploadSize   = 50 * 1024 * 1024 // 50 MB
    defaultDownload = 10 * 1024 * 1024 // 10 MB
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
    response := map[string]string{
        "message":   "pong",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
    }
    json.NewEncoder(w).Encode(response)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/octet-stream")

    size := defaultDownload
    if val := r.URL.Query().Get("size"); val != "" {
        if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 && parsed <= defaultDownload {
            size = parsed
        }
    }

    data := make([]byte, size)
    w.Write(data)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

    // Limit the upload size to 50MB
    r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
    start := time.Now()

    n, err := io.Copy(io.Discard, r.Body)
    duration := time.Since(start)

    if err != nil {
        http.Error(w, "Upload too large or failed to read upload", http.StatusRequestEntityTooLarge)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "received_bytes": n,
        "upload_time_ms": duration.Milliseconds(),
    })
}

func main() {
    http.HandleFunc("/ping", pingHandler)
    http.HandleFunc("/download", downloadHandler)
    http.HandleFunc("/upload", uploadHandler)

    port := os.Getenv("PORT")
    if port == "" {
        port = "10000"
    }

    http.ListenAndServe(":"+port, nil)
}