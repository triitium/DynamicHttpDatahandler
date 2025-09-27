package main

import (
    "log"
    "net/http"
    "os"
    "datahandler/config"
    "datahandler/handlers"
)

func main() {
    dbConn, err := config.ConnectDB()
    if err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }
    defer dbConn.Close()

    http.HandleFunc("/datastream/sensor", handlers.StandardDataHandler(dbConn))
    // http.HandleFunc("/api/fourier", handlers.StandardCalculationHandler(dbConn))

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }
    addr := ":" + port

    log.Printf("Spectrohub Go server running on %s\n", addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}
