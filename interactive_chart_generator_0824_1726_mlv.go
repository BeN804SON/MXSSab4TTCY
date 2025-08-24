// 代码生成时间: 2025-08-24 17:26:08
package main

import (
    "fmt"
    "net/http"
    "strings"
    "log"
    "encoding/json"
    "github.com/gorilla/mux"
)

// ChartData represents the data structure for chart input
type ChartData struct {
    Type    string `json:"type"`
    Data    [][]float64 `json:"data"`
    Options map[string]interface{} `json:"options"`
}

// ChartResponse represents the structure of the chart generation response
type ChartResponse struct {
    Message string `json:"message"`
    Chart   string `json:"chart"`
}

// GenerateChart generates an HTML chart using the provided data
func GenerateChart(data ChartData) (ChartResponse, error) {
    // Basic validation of chart type
    if data.Type == "" {
        return ChartResponse{}, fmt.Errorf("chart type is required")
    }
    
    // Create a simple HTML chart representation
    chartHTML := fmt.Sprintf("<html><body><canvas id='chart' width='400' height='400'></canvas><script>var ctx = document.getElementById('chart').getContext('2d');
var chart = new Chart(ctx, {
    type: '%s',
    data: {
        labels: [%s],
        datasets: [{
            label: 'Data',
            data: %s,
            backgroundColor: [
                'rgba(255, 99, 132, 0.2)',
                'rgba(54, 162, 235, 0.2)',
                'rgba(255, 206, 86, 0.2)'
            ],
            borderColor: [
                'rgba(255, 99, 132, 1)',
                'rgba(54, 162, 235, 1)',
                'rgba(255, 206, 86, 1)'
            ],
            borderWidth: 1
        }]
    },
    options: %s
});</script></body></html>", data.Type, strings.Join(strings.Split(strings.Repeat("', '", len(data.Data[0])-1), "'"), "
"), fmt.Sprintf("%v", data.Data), fmt.Sprintf("%v", data.Options))
    return ChartResponse{Message: "Chart generated successfully", Chart: chartHTML}, nil
}

// ChartHandler handles HTTP requests to generate charts
func ChartHandler(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var data ChartData
    err := decoder.Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()
    
    chartResponse, err := GenerateChart(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(chartResponse)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/chart", ChartHandler).Methods("POST")
    
    fmt.Println("Starting the interactive chart generator on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
