package config

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "encoding/json"
    "datahandler/db"
    "os"
)

type MetaData struct {
    APIKey   string          `json:"api_key"`
    Content  json.RawMessage `json:"content"`
}

type GeneralData struct {
    Name   string    `json:"name"`
    Values []float64 `json:"values"`
}

type BME280Data struct {
    Temperature float64 `json:"temperature"`
    Pressure   float64 `json:"pressure"`
    Humidity   float64 `json:"humidity"`
}

type BME688Data struct {
    Temperature     float64 `json:"temperature"`
    Pressure       float64 `json:"pressure"`
    Humidity       float64 `json:"humidity"`
    Gas_Resistance float64 `json:"gas_resistance"`
}

var gdStr = `
    CREATE TABLE IF NOT EXISTS %s (
        id SERIAL PRIMARY KEY,
        datapoints REAL[],
        timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
`

var TypeMapping = map[string]func() interface{}{
    "GeneralData": func() interface{} { return &GeneralData{} },
    "BME280Data":  func() interface{} { return &BME280Data{} },
    "BME688Data":  func() interface{} { return &BME688Data{} },
}

func ConnectDB() (*sql.DB, error) {
    dbURL := os.Getenv("POSTGRES_URL")
    dbUser := os.Getenv("POSTGRES_USER")
    dbPassword := os.Getenv("POSTGRES_PASSWORD")
    dbName := os.Getenv("POSTGRES_DB")

    connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbURL, dbName)
    return sql.Open("postgres", connStr)
}

type StandardMappingLogic struct {
    TableName       string
    DataType        string
    InsertFunc      func(db *sql.DB, table string, data interface{}) error
    CreateIfMissing bool
    TableString     *string
}

var StandardMapping = map[string]StandardMappingLogic{
    "64c1ecea-ba77-4271-acd9-7f06c6d1f004": {
        TableName:          "sensor.esp8266_bme688_envsense_m001",
        DataType:           "BME688Data",
        InsertFunc:         db.InsertStandardData,
        CreateIfMissing:    false,
        TableString:        nil,
    },
    "b012a270-b4f8-492c-b039-1bc7ef9f6000": {
        TableName:          "sensor.esp8266_bme280_sws1",
        DataType:           "BME280Data",
        InsertFunc:         db.InsertStandardData,
        CreateIfMissing:    false,
        TableString:        nil,
    },
    "92e11fa8-395e-4408-90c5-6368b3de1096": {
        TableName:          "%s_PLACEHOLDER",
        DataType:           "GeneralData",
        InsertFunc:         db.InsertStandardData,
        CreateIfMissing:    true,
        TableString:        &gdStr,
    },
}