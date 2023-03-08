package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	// "encoding/json"

	"github.com/almadigabor/maintenance-dash-go/internal/currentversions"
	"github.com/almadigabor/maintenance-dash-go/internal/data"
	"github.com/almadigabor/maintenance-dash-go/internal/latestversions"
	"github.com/almadigabor/maintenance-dash-go/internal/metrics"
	log "github.com/sirupsen/logrus"
)

var (
	cluster    *bool
	kubeconfig *string
)

const refreshIntervalInSeconds = 10 //5 * 60

func parseFlags() {
	// initialize flags
	cluster = flag.Bool("in-cluster", false, "Specify if the code is running inside a cluster or from outside.")
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
}

func syncAppsVersionInfo() {
	for {
		ctx := context.Background()
		var appsVersionInfo []*data.AppVersionInfo
		appCurrentInfos := currentversions.GetCurrentVersions(ctx, *cluster, *kubeconfig)

		for _, appCurrentInfo := range appCurrentInfos {
			appVersionInfo := latestversions.GetForApp(*appCurrentInfo)
			appsVersionInfo = append(appsVersionInfo, appVersionInfo)
		}

		metrics.UpdateMetrics(appsVersionInfo)
		time.Sleep(time.Duration(refreshIntervalInSeconds * float64(time.Second)))
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	parseFlags()

	go syncAppsVersionInfo()

	prometheusHandler := metrics.CreateMetrics()
	// setup metrics endpoint and start server
	http.Handle("/metrics", prometheusHandler)
	port := ":2112"
	log.Infof("Starting listening on %v", port)
	http.ListenAndServe(port, nil)
}
