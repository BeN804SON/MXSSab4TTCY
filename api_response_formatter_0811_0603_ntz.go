// 代码生成时间: 2025-08-11 06:03:41
package main

import (
    "net/http"
    "encoding/json"
    "log"
    "github.com/gorilla/mux"
)

// ApiResponse 结构体用于格式化API响应
type ApiResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// NewApiResponse 创建一个新的ApiResponse实例
func NewApiResponse(status, message string, data interface{}) *ApiResponse {
    return &ApiResponse{
        Status:  status,
        Message: message,
        Data:    data,
    }
}

// ErrorResponse 创建一个错误响应
func ErrorResponse(w http.ResponseWriter, err error, status int) {
    apiResponse := NewApiResponse("error", err.Error(), nil)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(apiResponse); err != nil {
        log.Printf("Failed to encode error response: %v", err)
    }
}

// OkResponse 创建一个成功的响应
func OkResponse(w http.ResponseWriter, data interface{}) {
    apiResponse := NewApiResponse("success", "Operation successful", data)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(apiResponse); err != nil {
        log.Printf("Failed to encode success response: %v", err)
    }
}

// HandleExampleEndpoint 处理示例端点
func HandleExampleEndpoint(w http.ResponseWriter, r *http.Request) {
    // 假设我们从数据库或其他服务中获取到了数据
    data := map[string]string{"key": "value"}
    OkResponse(w, data)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/example", HandleExampleEndpoint).Methods("GET")

    log.Println("Server starting on port 8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
