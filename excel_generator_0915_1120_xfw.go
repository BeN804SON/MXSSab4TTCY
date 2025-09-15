// 代码生成时间: 2025-09-15 11:20:06
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/360EntSecGroup-Skylar/excelize"
    "github.com/gorilla/mux"
)

// ExcelGenerator 用于生成Excel文件的结构体
type ExcelGenerator struct {
    // 在这里可以添加其他属性，例如模板路径等
}

// GenerateExcel 生成Excel文件
func (e *ExcelGenerator) GenerateExcel(w http.ResponseWriter, r *http.Request) {
    var data map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 使用excelize创建一个新的Excel文件
    f := excelize.NewFile()
    defer f.Close()

    // 设置Excel文件的标题
    f.SetSheetName(0, "Sheet1")

    // 根据传入的数据生成Excel内容
    for key, value := range data {
        // 假设key是列名，value是数据列表
        if list, ok := value.([]interface{}); ok {
            for i, item := range list {
                // 将数据写入Excel文件
                f.SetCellValue("Sheet1", fmt.Sprintf("A%d", 2+i), key)
                f.SetCellValue("Sheet1", fmt.Sprintf("B%d", 2+i), item)
            }
        }
    }

    // 将Excel文件写入响应
    w.Header().Set("Content-Disposition", "attachment; filename=generated_excel.xlsx")
    w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    if err := f.Write(w); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func main() {
    // 实例化ExcelGenerator
    excelGen := new(ExcelGenerator)

    // 设置路由和处理函数
    router := mux.NewRouter()
    router.HandleFunc("/generate_excel", excelGen.GenerateExcel).Methods("POST")

    // 启动服务器
    log.Println("Server starting on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
