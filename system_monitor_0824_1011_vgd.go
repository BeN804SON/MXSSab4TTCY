// 代码生成时间: 2025-08-24 10:11:34
package main

import (
    "fmt"
    "net/http"
    "os"
    "runtime"
    "time"
    "github.com/gorilla/mux"
)

// SystemInfo contains system performance data
type SystemInfo struct {
    Uptime    time.Duration `json:"uptime"`
    Memory    uint64        `json:"memory"`
    CPUUsage float64       `json:"cpu_usage"`
}

// getSystemInfo retrieves system performance data
func getSystemInfo() (SystemInfo, error) {
    uptime, err := os.Uptime()
    if err != nil {
        return SystemInfo{}, err
    }

    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    memory := memStats.Alloc

    // CPU usage calculation is platform-specific and complex;
    // for simplicity, we'll use a dummy value
    cpuUsage := 0.5 // 50% CPU usage

    return SystemInfo{
        Uptime:    uptime,
        Memory:    memory,
        CPUUsage: cpuUsage,
    }, nil
}

// systemInfoHandler handles HTTP requests for system information
func systemInfoHandler(w http.ResponseWriter, r *http.Request) {
    si, err := getSystemInfo()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(si)
}

func main() {
    router := mux.NewRouter()

    // Define the route for system info
    router.HandleFunc("/system", systemInfoHandler).Methods("GET")

    fmt.Println("Starting system performance monitor on :8080")
    err := http.ListenAndServe(":8080", router)
    if err != nil {
        fmt.Println("Failed to start server: ", err)
   }
}
