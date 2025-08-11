// 代码生成时间: 2025-08-11 20:42:45
package main

import (
    "fmt"
    "log"
    "net/http"
    "runtime"

    "github.com/gorilla/mux"
)

// MemoryUsageAnalyzer 结构体用于存储路由和内存分析功能
type MemoryUsageAnalyzer struct {
    // 可以添加其他属性
}

// NewMemoryUsageAnalyzer 创建并返回一个 MemoryUsageAnalyzer 实例
func NewMemoryUsageAnalyzer() *MemoryUsageAnalyzer {
    return &MemoryUsageAnalyzer{}
}

// MemoryStatsHandler 处理内存使用情况的HTTP请求
func (a *MemoryUsageAnalyzer) MemoryStatsHandler(w http.ResponseWriter, r *http.Request) {
    // 获取内存使用情况
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    // 构建响应数据
    response := struct {
        Alloc      uint64 `json:"alloc"`      // 从程序启动到现在分配的内存总量
        Sys        uint64 `json:"sys"`        // 从操作系统获得的内存总量
        HeapAlloc  uint64 `json:"heap_alloc"` // 堆内存分配量
        HeapSys    uint64 `json:"heap_sys"`    // 堆内存总量
        HeapIdle   uint64 `json:"heap_idle"`   // 堆内存空闲量
        HeapInuse  uint64 `json:"heap_inuse"`  // 堆内存使用量
        HeapReleased uint64 `json:"heap_released"` // 已释放的堆内存量
        // 可以添加更多内存使用情况的字段
    }{
        Alloc:      m.Alloc,
        Sys:        m.Sys,
        HeapAlloc:  m.HeapAlloc,
        HeapSys:    m.HeapSys,
        HeapIdle:   m.HeapIdle,
        HeapInuse:  m.HeapInuse,
        HeapReleased: m.HeapReleased,
    }

    // 将内存使用情况以JSON格式返回
    if err := encodeJSONResponse(w, response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// encodeJSONResponse 将数据编码为JSON并写入响应体
func encodeJSONResponse(w http.ResponseWriter, data interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(data)
}

func main() {
    router := mux.NewRouter()
    analyzer := NewMemoryUsageAnalyzer()

    // 定义内存使用情况分析的路由
    router.HandleFunc("/memory", analyzer.MemoryStatsHandler).Get()

    // 启动HTTP服务器
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
