// 代码生成时间: 2025-09-23 10:11:09
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"runtime"
	"log"

	"github.com/gorilla/mux"
)

// SystemStatus holds the data structure for system status
type SystemStatus struct {
	Uptime          time.Duration
	MemoryUsage    runtime.MemStats
	NumGoroutine    int
	NumCgoCall      int
	NumThread       int
	NumGoroutinePeak int
}	

// GetSystemStatus returns the current system status
func GetSystemStatus(w http.ResponseWriter, r *http.Request) {
	var status SystemStatus
	status.Uptime = time.Since(startTime)
	runtime.ReadMemStats(&status.MemoryUsage)
	status.NumGoroutine = runtime.NumGoroutine()
	status.NumCgoCall = runtime.NumCgoCall()
	status.NumThread = runtime.NumThread()
	status.NumGoroutinePeak = runtime.GOMAXPROCS(0) * (runtime.NumGoroutine() / runtime.GOMAXPROCS(0))

	// Marshal system status to JSON and write to response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}	

func main() {
	// Define the router
	r := mux.NewRouter()
	r.HandleFunc("/status", GetSystemStatus).Methods("GET")

	// Define a start time for uptime calculation
	var startTime time.Time
	startTime = time.Now()

	// Start the server
	log.Println("Starting system performance monitoring tool...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
