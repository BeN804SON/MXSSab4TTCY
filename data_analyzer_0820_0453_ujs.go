// 代码生成时间: 2025-08-20 04:53:20
package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
)

// 定义数据分析请求的结构体
type DataAnalysisRequest struct {
    Data []float64 `json:"data"`
}

// 定义数据分析结果的结构体
type AnalysisResult struct {
    Mean     float64 `json:"mean"`
    Median  float64 `json:"median"`
    Variance float64 `json:"variance"`
    StandardDeviation float64 `json:"standardDeviation"`
}

// 数据分析函数，计算平均值、中位数、方差和标准差
func analyzeData(data []float64) (AnalysisResult, error) {
    if len(data) == 0 {
        return AnalysisResult{}, fmt.Errorf("data slice is empty")
    }

    // 计算平均值
    var sum float64
    for _, v := range data {
        sum += v
    }
    mean := sum / float64(len(data))

    // 计算中位数
    sortedData := make([]float64, len(data))
    copy(sortedData, data)
    sort.Float64s(sortedData)
    median := sortedData[len(sortedData)/2]
    if len(sortedData)%2 == 0 {
        median = (sortedData[len(sortedData)/2-1] + sortedData[len(sortedData)/2]) / 2
    }

    // 计算方差
    var variance float64
    for _, v := range data {
        variance += (v - mean) * (v - mean)
    }
    variance /= float64(len(data) - 1)

    // 计算标准差
    stdDev := math.Sqrt(variance)

    return AnalysisResult{
        Mean:     mean,
        Median:  median,
        Variance: variance,
        StandardDeviation: stdDev,
    }, nil
}

// 数据分析器的HTTP处理函数
func dataAnalysisHandler(w http.ResponseWriter, r *http.Request) {
    var req DataAnalysisRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := analyzeData(req.Data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(result); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/analyze", dataAnalysisHandler).Methods("POST")

    log.Println("Starting data analysis server on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
