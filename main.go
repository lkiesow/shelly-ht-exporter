package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	humidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "humidity",
		Help: "The measured humidity",
	}, []string{"sensor"})
	temperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "temperature",
		Help: "The measured temperature",
	}, []string{"sensor"})
	countUpdated = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "count_updated",
		Help: "The number of updates to a sensor value",
	}, []string{"sensor"})
	lastUpdated = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "last_updated_time_seconds",
		Help: "Timestamp of the last update in seconds since epoch.",
	}, []string{"sensor"})
)

func init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(humidity)
	prometheus.MustRegister(temperature)
	prometheus.MustRegister(countUpdated)
	prometheus.MustRegister(lastUpdated)
}

func parseValue(w *http.ResponseWriter, strval string) (float64, error) {
	if val, err := strconv.ParseFloat(strval, 64); err == nil {
		return val, nil
	} else {
		(*w).WriteHeader(http.StatusBadRequest)
		(*w).Write([]byte("Unable to parse float value"))
		return 0, err
	}
}

func shelly(w http.ResponseWriter, req *http.Request) {
	var err error
	var hum, temp float64
	var id string

	fmt.Println("New sensor data:", req.URL.Query())

	// parse sensor values
	query := req.URL.Query()
	if hum, err = parseValue(&w, query.Get("hum")); err != nil {
		return
	}
	if temp, err = parseValue(&w, query.Get("temp")); err != nil {
		return
	}
	if id = query.Get("id"); id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("None of the query parameters `hum`, `temp` and `id` may be empty"))
		return
	}

	/// update metrics
	humidity.WithLabelValues(id).Set(hum)
	temperature.WithLabelValues(id).Set(temp)
	countUpdated.WithLabelValues(id).Inc()
	lastUpdated.WithLabelValues(id).Set(float64(time.Now().Unix()))
}

func main() {
	http.HandleFunc("/", shelly)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8090", nil)
}
