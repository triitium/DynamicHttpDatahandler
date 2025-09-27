package handlers

import (
    "log"
    "bytes"
    "datahandler/config"
    "datahandler/db"
    "database/sql"
    "encoding/json"
    "net/http"
)

func StandardDataHandler(conn *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var meta config.MetaData
        if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }

        cfg, ok := config.StandardMapping[meta.APIKey]
        if !ok {
            http.Error(w, "Invalid API key", http.StatusUnauthorized)
            return
        }

        mapper, ok := config.TypeMapping[cfg.DataType]
        if !ok {
            panic("unknown data type: " + cfg.DataType)
        }
        
        data := mapper() 
        if err := json.NewDecoder(bytes.NewReader(meta.Content)).Decode(&data); err != nil {
            log.Fatal(err)
        }

        if cfg.CreateIfMissing {
            if err := db.CreateTable(conn, cfg.TableName, cfg.TableString); err != nil {
                http.Error(w, "Failed to create table: "+err.Error(), http.StatusInternalServerError)
                return
            }
        }

        if err := cfg.InsertFunc(conn, cfg.TableName, data); err != nil {
            http.Error(w, "db insert Failed: "+err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Data processed"))
    }
}