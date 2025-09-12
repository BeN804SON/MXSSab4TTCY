// 代码生成时间: 2025-09-13 07:48:05
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"
)

// BackupData 结构体用于备份数据
type BackupData struct {
    Data string `json:"data"`
}

// RestoreData 结构体用于恢复数据
type RestoreData struct {
    FileName string `json:"fileName"`
}

// BackupHandler 处理器用于备份数据
func BackupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var backupData BackupData
    if err := json.NewDecoder(r.Body).Decode(&backupData); err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // 保存备份数据到文件
    fileName := fmt.Sprintf("backup_%s.json", backupData.Data)
    file, err := os.Create(fileName)
    if err != nil {
        http.Error(w, "Failed to create backup file", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    if _, err := file.WriteString(backupData.Data); err != nil {
        http.Error(w, "Failed to write backup data", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Backup successful: %s", fileName)
}

// RestoreHandler 处理器用于恢复数据
func RestoreHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var restoreData RestoreData
    if err := json.NewDecoder(r.Body).Decode(&restoreData); err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // 读取备份文件并恢复数据
    fileName := restoreData.FileName
    file, err := ioutil.ReadFile(fileName)
    if err != nil {
        http.Error(w, "Failed to read backup file", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Restore successful: %s", string(file))
}

func main() {
    // 设置路由
    router := http.NewServeMux()
    router.HandleFunc("/backup", BackupHandler)
    router.HandleFunc("/restore", RestoreHandler)

    // 启动服务器
    fmt.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
