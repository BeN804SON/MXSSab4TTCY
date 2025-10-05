// 代码生成时间: 2025-10-06 02:18:23
package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gorilla/mux"
)

// 创建一个简单的HTTP服务，用于端到端测试
func setupServer() *httptest.Server {
    router := mux.NewRouter()
    router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, world!")
    })
    return httptest.NewServer(router)
}

// 端到端测试函数
func TestEndToEnd(t *testing.T) {
    server := setupServer()
    defer server.Close()
    
    // 发送GET请求到服务器
    resp, err := http.Get(server.URL + "/test")
    if err != nil {
        t.Fatalf("Error making GET request: %v", err)
    }
    defer resp.Body.Close()
    
    // 检查HTTP响应状态码
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
    }
    
    // 读取响应体
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Fatalf("Error reading response body: %v", err)
    }
    expectedBody := "Hello, world!"
    if string(body) != expectedBody {
        t.Errorf("Expected body '%s', got '%s'", expectedBody, string(body))
    }
}

func main() {
    // 运行端到端测试
    testing.Main()
}