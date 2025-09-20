package db

import (
    "database/sql"
    "fmt"
    "strings"
)

func CreateTable(db *sql.DB, tableName string) error {
    query := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            id SERIAL PRIMARY KEY,
            datapoints REAL[],
            timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
    `, tableName)

    _, err := db.Exec(query)
    if err != nil {
        return fmt.Errorf("failed to create table %s: %w", tableName, err)
    }
    return nil
}

// Insert a single row into a table
func InsertSensorData(db *sql.DB, table string, values []float64) error {
    // Build query dynamically
    cols := []string{}
    placeholders := []string{}
    args := []interface{}{}
    for i, v := range values {
        colName := fmt.Sprintf("val%d", i+1)
        cols = append(cols, colName)
        placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
        args = append(args, v)
    }

    query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ","), strings.Join(placeholders, ","))
    _, err := db.Exec(query, args...)
    return err
}

func InsertSensorDataArray(db *sql.DB, table string, values []float64) error {
    // Convert []float64 to PostgreSQL array literal
    arrayLiteral := "{"
    for i, v := range values {
        if i > 0 {
            arrayLiteral += ","
        }
        arrayLiteral += fmt.Sprintf("%f", v)
    }
    arrayLiteral += "}"

    // Build and execute query
    query := fmt.Sprintf("INSERT INTO %s (datapoints) VALUES ($1)", table)
    _, err := db.Exec(query, arrayLiteral)
    return err
}
