package utils

import (
    "encoding/json"
    "net/http"
)

// RespondWithSuccess writes a success response to the http.ResponseWriter
func RespondWithSuccess(w http.ResponseWriter, message string, status int, data interface{}) {
    response := map[string]interface{}{
        "message":   message,
        "isSuccess": true,
        "status":    status,
        "data":      data,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}

// RespondWithError writes an error response to the http.ResponseWriter
func RespondWithError(w http.ResponseWriter, message string, status int) {
    response := map[string]interface{}{
        "message":   message,
        "isSuccess": false,
        "status":    status,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}
