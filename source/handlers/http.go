package handlers

import (
    "encoding/json"
    "net/http"
    "spectserver_source/config"
    "spectserver_source/db"
    "spectserver_source/fft"
)

type SensorData struct {
    Values []float64 `json:"values"`
    APIKey string    `json:"api_key"`
}

func SensorHandler(db **sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var data SensorData
        if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }

        // Validate API key
        cfg, ok := config.InsertMethods[data.APIKey]
        if !ok {
            http.Error(w, "Invalid API key", http.StatusUnauthorized)
            return
        }

        table := cfg.TableName

        if cfg.CreateIfMissing {
            if err := db.CreateTable(dbConn.DB, table); err != nil {
                http.Error(w, "Failed to create table: "+err.Error(), http.StatusInternalServerError)
                return
            }
        }

        if err := cfg.InsertFunc(dbConn.DB, table, data.Values); err != nil {
            http.Error(w, "DB insert Failed: "+err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Data processed"))
    }
}
