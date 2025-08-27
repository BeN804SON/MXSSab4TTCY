// 代码生成时间: 2025-08-28 04:49:21
package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/tealeg/xlsx/v3"
    "github.com/gorilla/mux"
)

// Define an error type to handle errors
type AppError struct {
    Message string
}

// ErrorResponse represents error response
type ErrorResponse struct {
    Error string `json:"error"`
}

// ExcelData represents the data to be written into the Excel file
type ExcelData struct {
    Data [][]string `json:"data"`
}

// writeExcelFile takes the data and writes it into an Excel file
func writeExcelFile(data ExcelData, filename string) error {
    file, err := xlsx.FilePaths(filename)
    if err != nil {
        return &AppError{Message: fmt.Sprintf("Error creating Excel file: %v", err)}
    }
    defer file.Close()

    sheet, err := file.AddSheet("Sheet1")
    if err != nil {
        return &AppError{Message: fmt.Sprintf("Error adding sheet: %v", err)}
    }

    for _, row := range data.Data {
        if err := sheet.AddRow(row...); err != nil {
            return &AppError{Message: fmt.Sprintf("Error adding row: %v", err)}
        }
    }

    return nil
}

// generateExcelHandler handles the request to generate an Excel file
func generateExcelHandler(w http.ResponseWriter, r *http.Request) {
    var data ExcelData
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    filename := fmt.Sprintf("%v.xlsx", time.Now().Format("20060102-150405"))
    if err := writeExcelFile(data, filename); err != nil {
        http.Error(w, err.(*AppError).Message, http.StatusInternalServerError)
        return
    }

    file, err := os.Open(filename)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error opening generated file: %v", err), http.StatusInternalServerError)
        return
    }
    defer file.Close()

    w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
    if _, err := io.Copy(w, file); err != nil {
        http.Error(w, fmt.Sprintf("Error writing file to response: %v", err), http.StatusInternalServerError)
        return
    }
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/generate", generateExcelHandler).Methods("POST")

    fmt.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
