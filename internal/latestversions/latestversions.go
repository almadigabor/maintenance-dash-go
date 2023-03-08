package latestversions

import (
	"context"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/almadigabor/maintenance-dash-go/internal/parseversion"
	"newreleases.io/newreleases"
)

type ProjectInfo struct {
	ID          string
	ReleaseName string
}

var (
	NewReleasesApiKey = os.Getenv("NEWREKEASES_API_KEY")
)

func GetAllProjects() []ProjectInfo {
	client := newreleases.NewClient(NewReleasesApiKey, nil)
	ctx := context.Background()
	var pi []ProjectInfo
	o := &newreleases.ProjectListOptions{
		Page: 1,
	}
	for {
		projects, lastPage, err := client.Projects.List(ctx, *o)
		if err != nil {
			log.Fatal(err)
		}

		for _, proj := range projects {
			releaseNameParts := strings.Split(proj.Name, "/")
			releaseName := releaseNameParts[len(releaseNameParts)-1]
			pi = append(pi, ProjectInfo{ID: proj.ID, ReleaseName: releaseName})
		}

		if o.Page >= lastPage {
			break
		}
		o.Page++
	}

	return pi
}

func GetSortedVersionsForApp(projectID string) []*semver.Version {
	client := newreleases.NewClient(NewReleasesApiKey, nil)
	ctx := context.Background()
	vs := []*semver.Version{}

	page := 1

	for {
		releases, lastPage, err := client.Releases.ListByProjectID(ctx, projectID, 1)
		if err != nil {
			log.Fatalln(err)
		}

		for _, release := range releases {
			parsedVersion := parseversion.ToSemver(release.Version)
			vs = append(vs, parsedVersion)
		}

		if page >= lastPage {
			break
		}
		page++
	}

	sort.Sort(sort.Reverse(semver.Collection(vs)))
	return vs
}
