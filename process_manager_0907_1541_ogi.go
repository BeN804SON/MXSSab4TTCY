// 代码生成时间: 2025-09-07 15:41:08
package main

import (
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "strings"
    "syscall"
    // Import the Gorilla mux package for routing
    "github.com/gorilla/mux"
)

// Process represents a system process
type Process struct {
    PID    int    `json:"pid"`
    Name   string `json:"name"`
    Status string `json:"status"`
}

// ProcessManager is a struct that contains necessary methods for managing processes
type ProcessManager struct {
    // Add fields if needed for more complex process management
}

// NewProcessManager creates a new instance of ProcessManager
func NewProcessManager() *ProcessManager {
    return &ProcessManager{}
}

// StartProcess starts a process with the given command
func (pm *ProcessManager) StartProcess(cmd string) (*Process, error) {
    // Split the command into parts
    parts := strings.Fields(cmd)
    if len(parts) == 0 {
        return nil, fmt.Errorf("invalid command")
    }

    // Start the process
    process, err := exec.Command(parts[0], parts[1:]...).Start()
    if err != nil {
        return nil, fmt.Errorf("failed to start process: %w", err)
    }

    // Retrieve the process ID
    pid := process.Process.Pid

    // Create a new Process struct and return it
    return &Process{PID: pid, Name: parts[0], Status: "running"}, nil
}

// StopProcess stops a process with the given PID
func (pm *ProcessManager) StopProcess(pid int) error {
    // Use syscall to send a SIGTERM to the process
    process, err := os.FindProcess(pid)
    if err != nil {
        return fmt.Errorf("failed to find process with PID %d: %w", pid, err)
    }
    if err := process.Signal(syscall.SIGTERM); err != nil {
        return fmt.Errorf("failed to stop process with PID %d: %w", pid, err)
    }

    return nil
}

// ListProcesses lists all running processes
func (pm *ProcessManager) ListProcesses() ([]Process, error) {
    // This is a simplified example and does not list all processes.
    // In a real-world scenario, you would need to interact with the system's process list.
    // For demonstration purposes, we'll return an empty list.
    return []Process{}, nil
}

// SetupRoutes sets up the HTTP routes for the process manager
func SetupRoutes(pm *ProcessManager) *mux.Router {
    router := mux.NewRouter()

    // Start a process
    router.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        cmd := vars["cmd"]
        process, err := pm.StartProcess(cmd)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Return the started process as JSON
        fmt.Fprintf(w, "{"pid":%d, "name":"%s", "status":"running"}", process.PID, process.Name)
    }).Methods("POST")

    // Stop a process
    router.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        pidStr := vars["pid"]
        pid, err := strconv.Atoi(pidStr)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if err := pm.StopProcess(pid); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        fmt.Fprintln(w, "Process stopped")
    }).Methods("POST")

    // List processes
    router.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
        processes, err := pm.ListProcesses()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Convert the processes to JSON and return them
        fmt.Fprintf(w, "{"processes":%s}", toJSON(processes))
    }).Methods("GET")

    return router
}

// toJSON converts a slice of Process structs to a JSON string
func toJSON(processes []Process) string {
    result, err := json.Marshal(processes)
    if err != nil {
        log.Fatalf("Failed to marshal processes to JSON: %v", err)
    }
    return string(result)
}

func main() {
    pm := NewProcessManager()
    router := SetupRoutes(pm)

    // Start the HTTP server
    log.Printf("Starting process manager on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
