package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	// Set up logging configuration
	logFile, err := os.Create("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	// Start logging system metrics in a goroutine
	go logSystemMetrics(logger)

	http.HandleFunc("/", indexHandler)
	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", http.FileServer(http.Dir("static"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/get_metrics", getMetricsHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func logSystemMetrics(logger *log.Logger) {
	for {
		cpuMetrics, err := cpu.Percent(time.Second, false)
		if err != nil {
			logger.Println(err)
		}

		memMetrics, err := mem.VirtualMemory()
		if err != nil {
			logger.Println(err)
		}

		diskMetrics, err := disk.Usage("/")
		if err != nil {
			logger.Println(err)
		}

		logger.Printf("CPU Usage: %.2f%%\n", cpuMetrics[0])
		logger.Printf("Memory Usage: %.2f%%\n", memMetrics.UsedPercent)
		logger.Printf("Disk Usage: %.2f%%\n", diskMetrics.UsedPercent)

		if cpuMetrics[0] > 80 || memMetrics.UsedPercent > 80 || diskMetrics.UsedPercent > 80 {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			message := fmt.Sprintf("Scale up your system: Device usage is more than 80%% at %s.", timestamp)
			logger.Println(message)
		}

		time.Sleep(2 * time.Second)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func getMetricsHandler(w http.ResponseWriter, r *http.Request) {
	cpuMetrics, err := cpu.Percent(time.Second, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	memMetrics, err := mem.VirtualMemory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	diskMetrics, err := disk.Usage("/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := ""
	if cpuMetrics[0] > 80 || memMetrics.UsedPercent > 80 || diskMetrics.UsedPercent > 80 {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		message = fmt.Sprintf("Scale up your system: Device usage is more than 80%% at %s.", timestamp)
	}

	data := map[string]interface{}{
		"cpu_metric":  cpuMetrics[0],
		"mem_metric":  memMetrics.UsedPercent,
		"disk_metric": diskMetrics.UsedPercent,
		"message":     message,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
