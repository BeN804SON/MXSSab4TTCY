// 代码生成时间: 2025-09-16 18:09:23
package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// StatisticalData defines the structure for statistical data
type StatisticalData struct {
	Data []float64
}

// DataAnalyticsService is a struct that will hold the data and methods for statistical analysis
type DataAnalyticsService struct {
	Data StatisticalData
}

// NewDataAnalyticsService initializes a new instance of DataAnalyticsService with sample data
func NewDataAnalyticsService() *DataAnalyticsService {
	// Sample data for demonstration purposes
	sampleData := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
	return &DataAnalyticsService{
		Data: StatisticalData{Data: sampleData},
	}
}

// CalculateMean calculates the mean of the data set
func (das *DataAnalyticsService) CalculateMean() (float64, error) {
	if len(das.Data.Data) == 0 {
		return 0, fmt.Errorf("data set is empty")
	}
	sum := 0.0
	for _, value := range das.Data.Data {
		sum += value
	}
	return sum / float64(len(das.Data.Data)), nil
}

// CalculateMedian calculates the median of the data set
func (das *DataAnalyticsService) CalculateMedian() (float64, error) {
	dataLen := len(das.Data.Data)
	if dataLen == 0 {
		return 0, fmt.Errorf("data set is empty")
	}
	sortedData := make([]float64, dataLen)
	copy(sortedData, das.Data.Data)
	sort.Float64s(sortedData)
	mid := dataLen / 2
	if dataLen%2 == 0 {
		return (sortedData[mid-1] + sortedData[mid]) / 2, nil
	}
	return sortedData[mid], nil
}

// CalculateMode calculates the mode of the data set
func (das *DataAnalyticsService) CalculateMode() (float64, error) {
	if len(das.Data.Data) == 0 {
		return 0, fmt.Errorf("data set is empty")
	}
	frequencyMap := make(map[float64]int)
	for _, value := range das.Data.Data {
		frequencyMap[value]++
	}
	maxFrequency := 0
	mode := 0.0
	for value, frequency := range frequencyMap {
		if frequency > maxFrequency {
			maxFrequency = frequency
			mode = value
		}
	}
	return mode, nil
}

// HandleStatisticsRequest handles HTTP requests to calculate statistics
func HandleStatisticsRequest(w http.ResponseWriter, r *http.Request) {
	service := NewDataAnalyticsService()

	vars := mux.Vars(r)
	statType := vars["type"]

	switch statType {
	case "mean":
		mean, err := service.CalculateMean()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatFloat(mean, 'f', 2, 64)))
	case "median":
		median, err := service.CalculateMedian()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatFloat(median, 'f', 2, 64)))
	case "mode":
		mode, err := service.CalculateMode()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatFloat(mode, 'f', 2, 64)))
	default:
		http.Error(w, "invalid statistic type", http.StatusBadRequest)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/statistics/{type}", HandleStatisticsRequest).Methods("GET")

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", router)
}
