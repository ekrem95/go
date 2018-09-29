package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("port", "8080", "The port to listen on for HTTP requests.")
	m    runtime.MemStats
	list [][]int
)

var (
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})
	memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage_mb",
		Help: "Current memory usage.",
	})
)

func init() {
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(memoryUsage)
}

func main() {
	flag.Parse()

	go func() {
		for {
			cpuTemp.Set(rand.Float64()*(70-45) + 45)
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		for {
			// https://golang.org/pkg/runtime/#MemStats
			for i := 0; i < 4; i++ {
				a := make([]int, 0, 99999)
				list = append(list, a)
			}

			runtime.ReadMemStats(&m)
			memoryUsage.Set(bToMb(m.Sys) + float64(rand.Intn(10)))

			// Force GC to clear up
			if bToMb(m.Alloc) >= 100 {
				list = nil
				runtime.GC()

				runtime.ReadMemStats(&m)
				memoryUsage.Set(bToMb(m.Sys))
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+*addr, nil))
}

func bToMb(b uint64) float64 {
	return float64(b / 1024 / 1024)
}
