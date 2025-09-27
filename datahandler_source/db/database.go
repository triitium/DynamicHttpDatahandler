package db

import (
    "time"
    "database/sql"
    "fmt"
    "reflect"
    "strings"
)

func CreateTable(conn *sql.DB, tableName string, tableString *string) error {
    now := time.Now()
    tableName = fmt.Sprintf(tableName, now.Format("2006_01_02_15_04_05"))
    query := fmt.Sprintf(*tableString, tableName)

    _, err := conn.Exec(query)
    if err != nil {
        return fmt.Errorf("failed to create table %s: %w", tableName, err)
    }
    return nil
}

func InsertStandardData(conn *sql.DB, table string, data interface{}) error {
    v := reflect.ValueOf(data)
    if v.Kind() != reflect.Struct {
        return fmt.Errorf("data must be a struct")
    }

    cols := []string{}
    placeholders := []string{}
    args := []interface{}{}

    for i := 0; i < v.NumField(); i++ {
        field := v.Type().Field(i)
        value := v.Field(i).Interface()

        colName := field.Tag.Get("conn")
        if colName == "" {
            colName = strings.ToLower(field.Name)
        }

        cols = append(cols, colName)
        placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)+1))
        args = append(args, value)
    }

    if len(cols) == 0 {
        return fmt.Errorf("no columns to insert")
    }

    query := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s)",
        table,
        strings.Join(cols, ","),
        strings.Join(placeholders, ","),
    )

    _, err := conn.Exec(query, args...)
    return err
}

func InsertSensorDataArray(conn *sql.DB, table string, values []float64) error {
    arrayLiteral := "{"
    for i, v := range values {
        if i > 0 {
            arrayLiteral += ","
        }
        arrayLiteral += fmt.Sprintf("%f", v)
    }
    arrayLiteral += "}"

    query := fmt.Sprintf("INSERT INTO %s (datapoints) VALUES ($1)", table)
    _, err := conn.Exec(query, arrayLiteral)
    return err
}