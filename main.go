package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type sensor struct {
	count uint
	hum   float64
	temp  float64
}

var sensors = make(map[string]sensor)

func metrics(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
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
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}

	fmt.Println("GET params were:", req.URL.Query())

	// if only one expected
	hum := req.URL.Query().Get("hum")
	temp := req.URL.Query().Get("temp")
	id := req.URL.Query().Get("id")
	if hum == "" || temp == "" || id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("None of the query parameters `hum`, `temp` and `id` may be empty"))
		return
	}

	s, ok := sensors[id]
	s.count += 1
	var err error
	if s.hum, err = parseValue(&w, temp); err != nil {
		return
	}
	if s.temp, err = parseValue(&w, temp); err != nil {
		return
	}
	sensors[id] = s
	fmt.Println("in map:", ok)
	fmt.Println("sensor:", s)
	fmt.Println("sensors:", sensors)
}

func main() {
	http.HandleFunc("/", shelly)
	http.HandleFunc("/metrics", metrics)
	http.ListenAndServe(":8090", nil)
}
