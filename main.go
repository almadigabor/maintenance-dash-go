package main

import (
	"flag"
	"net/http"

	// "encoding/json"

	"github.com/Masterminds/semver/v3"
	"github.com/almadigabor/maintenance-dash-go/pkg/data"
	"github.com/almadigabor/maintenance-dash-go/pkg/latestversions"
	"github.com/almadigabor/maintenance-dash-go/pkg/metrics"
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
		panic(err.Error())
	}

	appsVersionInfo = latestversions.GetAppLatestVersions(c)

	//clientSet := currentversions.GetClientSet(*cluster, *kubeconfig)
	//nodes, _ := clientSet.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	//
	//for _, node := range nodes.Items {
	//	fmt.Printf("%s: %s\n", node.Name, node.Status.NodeInfo.KubeletVersion)
	//}
	//
	//appsVersionInfo := make([]data.AppVersionInfo, 0)
	prometheusHandler := metrics.CreateAppsVersionMetrics(appsVersionInfo)
	// setup metrics endpoint and start server
	http.Handle("/metrics", prometheusHandler)
	port := ":2112"
	log.Infof("Starting listening on %v", port)
	http.ListenAndServe(port, nil)
}
