package main

import (
	"flag"
	"fmt"
	"net/http"

	// "encoding/json"

	"github.com/Masterminds/semver/v3"
	"github.com/almadigabor/maintenance-dash-go/internal/data"
	"github.com/almadigabor/maintenance-dash-go/internal/latestversions"
	"github.com/almadigabor/maintenance-dash-go/internal/metrics"
	log "github.com/sirupsen/logrus"
)

var (
	cluster         *bool
	kubeconfig      *string
	appsVersionInfo []*data.AppVersionInfo = make([]*data.AppVersionInfo, 0)
)

type Versions []semver.Version

func parseFlags() {
	// initialize flags
	cluster = flag.Bool("in-cluster", false, "Specify if the code is running inside a cluster or from outside.")
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	parseFlags()

	c, err := data.ReadConf("config.yaml")
	if err != nil {
		log.Panic(err.Error())
	}

	//ctx := context.Background()
	//clientSet := currentversions.NewClientSet(*cluster, *kubeconfig)

	projectInfos := latestversions.GetAllProjects()
	for _, pi := range projectInfos {
		fmt.Println(pi.ReleaseName)
	}
	//currentversions.GetSvcsToScan(ctx, clientSet)
	//currentversions.AddK8sNodeVersionInfo(ctx, clientSet, appsVersionInfo)

	appsVersionInfo = append(appsVersionInfo, latestversions.GetAppLatestVersions(c)...)

	prometheusHandler := metrics.CreateAppsVersionMetrics(appsVersionInfo)
	// setup metrics endpoint and start server
	http.Handle("/metrics", prometheusHandler)
	port := ":2112"
	log.Infof("Starting listening on %v", port)
	http.ListenAndServe(port, nil)
}
