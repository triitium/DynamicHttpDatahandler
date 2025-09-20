package config

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "os"
)

func ConnectDB() (*sql.DB, error) {
    dbURL := os.Getenv("POSTGRES_URL")
    dbUser := os.Getenv("POSTGRES_USER")
    dbPassword := os.Getenv("POSTGRES_PASSWORD")
    dbName := os.Getenv("POSTGRES_DB")

    connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbURL, dbName)
    return sql.Open("postgres", connStr)
}

type InsertConfig struct {
    TableName       string
    InsertFunc      func(db *sql.DB, table string, values []float64) error
    CreateIfMissing bool
}

var InsertMethods = map[string]InsertConfig{
    "EuCcPAP5QPnC4HJWc08u": {
        TableName:      "esp32_athmos_spectro_001",
        InsertFunc:     db.InsertSensorDataArray,
        CreateNew:      false,
    },
    "ImO4W0xrtN5mYp9QvjGq": {
        TableName:      "sensor2_table",
        InsertFunc:     db.InsertSensorDataArray,
        CreateNew:      false,
    },
}
