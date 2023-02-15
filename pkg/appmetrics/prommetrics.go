package appmetrics

import (
	"net/http"

	"github.com/almadigabor/maintenance-dash-go/pkg/data"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func CreateAppsVersionMetrics(appsVersionInfo []data.AppVersionInfo) http.Handler {
	// get rid of the default metrics
	r := prometheus.NewRegistry()
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	// add the metrics for all applications with Gauge type
	addAppMetrics(r, appsVersionInfo)
	return handler
}

func addAppMetrics(collector *prometheus.Registry, appsVersionInfo []data.AppVersionInfo) {
	for _, appVersionInfo := range appsVersionInfo {
		var latestMajorVersion, latestMinorVersion, latestPatchVersion string

		// Convert the semantic versions to string
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
