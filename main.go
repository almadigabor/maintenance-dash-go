package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	// "encoding/json"
	"os"

	"github.com/Masterminds/semver/v3"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	AppsVersionInfo []AppVersionInfo
)

type AppVersionInfo struct {
	AppName            string
	CurrentVersion     *semver.Version
	LatestMajorVersion *semver.Version
	LatestMinorVersion *semver.Version
	LatestPatchVersion *semver.Version
}

type Release struct {
	Version      string   `json:"version,omitempty"`
	Date         string   `json:"date,omitempty"`
	IsPrerelease bool     `json:"is_prerelease,omitempty"`
	HasNote      bool     `json:"has_note,omitempty"`
	Cve          []string `json:"cve,omitempty"`
}

type NewReleasesResponse struct {
	Releases []Release `json:"releases,omitempty"`
}

type Versions []semver.Version

func addAppMetrics(collector *prometheus.Registry) {
	for _, appVersionInfo := range AppsVersionInfo {
		var latestMajorVersion, latestMinorVersion, latestPatchVersion string
		if appVersionInfo.LatestMajorVersion != nil {
			latestMajorVersion = appVersionInfo.LatestMajorVersion.String()
		}
		if appVersionInfo.LatestMinorVersion != nil {
			latestMinorVersion = appVersionInfo.LatestMinorVersion.String()
		}
		if appVersionInfo.LatestPatchVersion != nil {
			latestPatchVersion = appVersionInfo.LatestPatchVersion.String()
		}
		var appMetric = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "app_version_info",
			Help: "Information about an application containing the current, latest major, latest minor, latest patch versions.",
			ConstLabels: prometheus.Labels{
				"appName":            appVersionInfo.AppName,
				"currentVersion":     appVersionInfo.CurrentVersion.String(),
				"latestMajorVersion": latestMajorVersion,
				"latestMinorVersion": latestMinorVersion,
				"latestPatchVersion": latestPatchVersion,
			},
		})

		collector.MustRegister(appMetric)
	}
}

// Returns the latest major, minor and patch versions for the application based on the current version from newreleases.io
func getLatestVersionsForApp(provider string, project string, repository string, currentVersion *semver.Version) AppVersionInfo {
	NewReleasesApiKey := os.Args[1]
	client := http.Client{}
	url := fmt.Sprintf("https://api.newreleases.io/v1/projects/%v/%v/%v/releases", provider, project, repository)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"X-Key":        {NewReleasesApiKey},
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(res.Body)
	var response = &NewReleasesResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		panic(err)
	}
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	vs := []*semver.Version{}
	for _, r := range response.Releases {
		v, err := semver.NewVersion(r.Version)
		if err != nil {
			fmt.Printf("Error parsing version: %s\n", err)
		} else if v.Prerelease() == "" {
			vs = append(vs, v)
		}
	}

	sort.Sort(sort.Reverse(semver.Collection(vs)))
	latestMajorVersion := vs[0]
	var latestMinorVersion *semver.Version
	var latestPatchVersion *semver.Version

	for i := 0; i < len(vs); i++ {
		if latestMinorVersion == nil && vs[i].Major() == currentVersion.Major() {
			latestMinorVersion = vs[i]
		}

		if latestPatchVersion == nil && vs[i].Major() == currentVersion.Major() && vs[i].Minor() == currentVersion.Minor() {
			latestPatchVersion = vs[i]
		}
	}

	return AppVersionInfo{
		AppName:            repository,
		CurrentVersion:     currentVersion,
		LatestMajorVersion: latestMajorVersion,
		LatestMinorVersion: latestMinorVersion,
		LatestPatchVersion: latestPatchVersion,
	}
}

func main() {
	currentKubernetesVersion, _ := semver.NewVersion("1.24.2")
	var kubernetesLatestVersions = getLatestVersionsForApp("github", "kubernetes", "kubernetes", currentKubernetesVersion)
	AppsVersionInfo = append(AppsVersionInfo, kubernetesLatestVersions)

	currentArgocdVersion, _ := semver.NewVersion("2.3.12")
	var argocdLatestVersions = getLatestVersionsForApp("github", "argoproj", "argo-cd", currentArgocdVersion)
	AppsVersionInfo = append(AppsVersionInfo, argocdLatestVersions)

	currentCertManagerVersion, _ := semver.NewVersion("1.10.0")
	var certManagerLatestVersions = getLatestVersionsForApp("github", "cert-manager", "cert-manager", currentCertManagerVersion)
	AppsVersionInfo = append(AppsVersionInfo, certManagerLatestVersions)

	// get rid of the default metrics
	r := prometheus.NewRegistry()
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	// add the metrics for all applications with Gauge type
	addAppMetrics(r)

	// setup metrics endpoint and start server
	http.Handle("/metrics", handler)
	http.ListenAndServe(":2112", nil)
}
