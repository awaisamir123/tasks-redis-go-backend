package main

import (
    "log"
    "net/http"
    "os"
    "task-management/routes"
    "github.com/joho/godotenv"

)


func main() {
    // Load environment variables from the .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    serverPort := os.Getenv("SERVER_PORT")

    // Initialize the router and register routes
    r := routes.SetupRouter()

    log.Printf("Server starting on port %s", serverPort)
    
    // Use '=' here instead of ':=' since 'err' is already declared
    err = http.ListenAndServe(":"+serverPort, r)
    if err != nil {
        log.Fatalf("Error starting server: %s", err)
    }
}