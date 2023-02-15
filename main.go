package main

import (
	"flag"

	// "encoding/json"

	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"
)

var (
	cluster    *bool
	kubeconfig *string
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

	c, err := readConf("config.yaml")
	if err != nil {
		panic(err.Error())
	}

	listPublicRepos(c)

	//clientSet := currentversions.GetClientSet(*cluster, *kubeconfig)
	//nodes, _ := clientSet.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	//
	//for _, node := range nodes.Items {
	//	fmt.Printf("%s: %s\n", node.Name, node.Status.NodeInfo.KubeletVersion)
	//}
	//
	//appsVersionInfo := make([]data.AppVersionInfo, 0)
	//prometheusHandler := appmetrics.CreateAppsVersionMetrics(appsVersionInfo)
	//// setup metrics endpoint and start server
	//http.Handle("/metrics", prometheusHandler)
	//http.ListenAndServe(":2112", nil)
}
