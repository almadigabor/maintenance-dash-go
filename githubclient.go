package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"os"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var (
	GitHubAccessToken = os.Getenv("GH_ACCESS_TOKEN")
)

func listPublicRepos(c *ReleaseConfig) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GitHubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	for _, configItem := range c.Apps {
		repo, _, err := client.Repositories.Get(ctx, configItem.OrgName, configItem.RepoName)
		log.Infof("Getting releases for repo: %v", *repo.Name)

		if err != nil {
			log.Panic(err.Error())
		}

		vs := []*semver.Version{}

		for i := 1; i < 10; i++ {
			opt := &github.ListOptions{
				PerPage: 10,
				Page:    i,
			}
			releases, _, err := client.Repositories.ListReleases(ctx, configItem.OrgName, *repo.Name, opt)
			if err != nil {
				log.Panic(err.Error())
			}

			for _, release := range releases {
				if !strings.HasPrefix(*release.Name, configItem.ReleasePrefix) {
					log.Infof("Skipping release with wrong prefix: %s", *release.Name)
					continue
				}
				versionString := strings.TrimPrefix(*release.Name, configItem.ReleasePrefix+"-")
				v, err := semver.NewVersion(versionString)
				if err != nil {
					log.Infof("Error parsing %s version: %s", versionString, err)
				} else if v.Prerelease() == "" {
					vs = append(vs, v)
				}
			}

			if len(vs) >= 10 {
				break
			}
		}

		log.Infof("Found %v semver releases with prefix %v for repo %v", len(vs), configItem.ReleasePrefix, *repo.Name)
		sort.Sort(sort.Reverse(semver.Collection(vs)))
		for _, ver := range vs {
			fmt.Printf("\t%s\n", ver.String())
		}
		//latestMajorVersion := vs[0]
		//var latestMinorVersion *semver.Version
		//var latestPatchVersion *semver.Version
		//
		//for i := 0; i < len(vs); i++ {
		//	if latestMinorVersion == nil && vs[i].Major() == currentVersion.Major() {
		//		latestMinorVersion = vs[i]
		//	}
		//
		//	if latestPatchVersion == nil && vs[i].Major() == currentVersion.Major() && vs[i].Minor() == currentVersion.Minor() {
		//		latestPatchVersion = vs[i]
		//	}
		//}
	}
}
