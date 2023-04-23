package main

import (
	"flag"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "yr"

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	addr      = flag.String("addr", ":9367", "The address to listen on for HTTP requests.")
	locations arrayFlags
)

func init() {
	flag.Var(&locations, "location", "lat,long,name")
}

func main() {
	flag.Parse()

	logger, _ := zap.NewProduction()
	logger.Info("Starting yr_exporter")

	var parsedLocations []location
	for _, l := range locations {
		p := strings.Split(l, ",")
		if len(p) != 3 {
			panic("expected exactly 3 parts to location (lat,long,name)")
		}
		parsedLocations = append(parsedLocations, location{
			lat:  p[0],
			long: p[1],
			name: p[2],
		})
	}

	prometheus.MustRegister(NewYrCollector(namespace, logger, parsedLocations))

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>YR Exporter</title></head>
            <body>
            <h1>YR Exporter</h1>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})
	srv := &http.Server{
		Addr:         *addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Listening on", zap.Stringp("addr", addr))
	logger.Fatal("failed to start server", zap.Error(srv.ListenAndServe()))
}
