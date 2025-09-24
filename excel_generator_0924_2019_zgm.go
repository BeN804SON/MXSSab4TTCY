// 代码生成时间: 2025-09-24 20:19:18
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/tealeg/xlsx"
    "github.com/gorilla/mux"
)

// ExcelGenerator 结构体用于存储生成的Excel文件信息
type ExcelGenerator struct {
    // 文件名和路径
    FileName string
    // Excel文件
    File *xlsx.File
}

// NewExcelGenerator 创建一个新的ExcelGenerator实例
func NewExcelGenerator() *ExcelGenerator {
    return &ExcelGenerator{
        FileName: time.Now().Format("2006-01-02T15-04-05") + ".xlsx",
    }
}

// GenerateExcel 生成Excel文件
func (e *ExcelGenerator) GenerateExcel() error {
    e.File = xlsx.NewFile()
    sheet, err := e.File.AddSheet("Sheet1")
    if err != nil {
        return err
    }
    // 向Excel添加数据
    // 示例：添加一行数据
    sheet.AddRow()
    sheet.AddCell("My Header")
    // 这里可以添加更多的行和单元格
    // ...
    // 写入文件
    f, err := os.Create(e.FileName)
    if err != nil {
        return err
    }
    defer f.Close()
    if err := e.File.Write(f); err != nil {
        return err
    }
    return nil
}

// ServeExcel 提供一个HTTP服务端点，用于下载生成的Excel文件
func (e *ExcelGenerator) ServeExcel(w http.ResponseWriter, r *http.Request) {
    err := e.GenerateExcel()
    if err != nil {
        log.Printf("Error generating Excel: %s", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.ServeFile(w, r, e.FileName)
    os.Remove(e.FileName) // 可选：删除临时文件
}

func main() {
    r := mux.NewRouter()
    eg := NewExcelGenerator()
    r.HandleFunc("/download", eg.ServeExcel) // 设置路由
    port := ":8080"
    fmt.Printf("Server starting on port %s...
", port)
    log.Fatal(http.ListenAndServe(port, r))
}