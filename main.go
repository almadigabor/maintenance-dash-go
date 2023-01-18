package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	// "encoding/json"
	"os"

	"github.com/Masterminds/semver/v3"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Release struct {
	Version      string   `json:"version,omitempty"`
	Date         string   `json:"date,omitempty"`
	IsPrerelease bool     `json:"is_prerelease,omitempty"`
	HasNote      bool     `json:"has_note,omitempty"`
	Cve          []string `json:"cve,omitempty"`
}

type Response struct {
	Releases []Release `json:"releases,omitempty"`
}

type Versions []semver.Version

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func makeApiCall() {
	NewReleasesApiKey := os.Args[1]
	client := http.Client{}
	provider := "github"
	project := "argoproj"
	repo := "argo-cd"

	url := fmt.Sprintf("https://api.newreleases.io/v1/projects/%v/%v/%v/releases", provider, project, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// req.Header.Set("X-Key", NewReleasesApiKey)

	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"X-Key":        {NewReleasesApiKey},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(res.Body)

	var response = &Response{}

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

	var versions []string
	for _, release := range response.Releases {
		versions = append(versions, release.Version)
	}

	vs := []*semver.Version{}
	for _, r := range versions {
		v, err := semver.NewVersion(r)
		if err != nil {
			fmt.Printf("Error parsing version: %s\n", err)
		} else if v.Prerelease() == "" {
			vs = append(vs, v)
		}
	}

	sort.Sort(semver.Collection(vs))
	latestVersion := vs[len(vs)-1]
	fmt.Println(latestVersion)
}

func main() {
	// example newreleases api call
	makeApiCall()

	// example metrics
	recordMetrics()

	// get rid of the default metrics
	r := prometheus.NewRegistry()
	r.MustRegister(opsProcessed)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	// setup metrics endpoint and start server
	http.Handle("/metrics", handler)
	http.ListenAndServe(":2112", nil)
}
