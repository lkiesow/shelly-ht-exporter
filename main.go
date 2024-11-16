// Shelly H&T Prometheus Exporter
// Copyright 2024 Lars Kiesow <lkiesow@uos.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"net/http"
	"os"
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
	// Get the listen address
	addr := os.Getenv("SHELLY_HT_EXPORTER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8090"
	}

	// register handlers
	http.HandleFunc("/", shelly)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Printf("Listening on on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
