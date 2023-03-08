package latestversions

import (
	"os"

	"github.com/Masterminds/semver/v3"
)

var (
	GitHubAccessToken = os.Getenv("GH_ACCESS_TOKEN")
)

type latestVersions struct {
	majorVersion *semver.Version
	minorVersion *semver.Version
	patchVersion *semver.Version
}

// func GetAppLatestVersions(c *data.ReleaseConfig) []*data.AppVersionInfo {
// 	ctx := context.Background()
// 	ts := oauth2.StaticTokenSource(
// 		&oauth2.Token{AccessToken: GitHubAccessToken},
// 	)
// 	tc := oauth2.NewClient(ctx, ts)
// 	client := github.NewClient(tc)

// 	var appsVersionInfo []*data.AppVersionInfo
// 	for _, configItem := range c.Apps {
// 		currentVersion, _ := semver.NewVersion("1.0.0")
// 		vs := findLatestVersions(currentVersion, getReposReleases(ctx, client, &configItem))
// 		appName := configItem.ReleasePrefix
// 		if len(appName) == 0 {
// 			appName = configItem.RepoName
// 		}
// 		appVersionInfo := data.AppVersionInfo{
// 			AppName:            appName,
// 			CurrentVersion:     currentVersion,
// 			LatestMajorVersion: vs.majorVersion,
// 			LatestMinorVersion: vs.minorVersion,
// 			LatestPatchVersion: vs.patchVersion,
// 		}
// 		appsVersionInfo = append(appsVersionInfo, &appVersionInfo)
// 	}

// 	return appsVersionInfo
// }

// func getReposReleases(ctx context.Context, client *github.Client, configItem *data.ReleaseConfigItem) []*semver.Version {
// 	repo, _, err := client.Repositories.Get(ctx, configItem.OrgName, configItem.RepoName)
// 	log.Infof("Getting releases for repo: %v", *repo.Name)

// 	if err != nil {
// 		log.Panic(err.Error())
// 	}

// 	vs := []*semver.Version{}

// 	// read 10 pages of 10 releases until 10 proper releases are found
// 	for i := 1; i < 10; i++ {
// 		opt := &github.ListOptions{
// 			PerPage: 10,
// 			Page:    i,
// 		}
// 		releases, _, err := client.Repositories.ListReleases(ctx, configItem.OrgName, *repo.Name, opt)
// 		if err != nil {
// 			log.Panic(err.Error())
// 		}

// 		// check release prefix and parse filtered releases to semver
// 		for _, release := range releases {
// 			if !strings.HasPrefix(*release.Name, configItem.ReleasePrefix) {
// 				log.Infof("Skipping release with wrong prefix: %s", *release.Name)
// 				continue
// 			}
// 			versionString := strings.TrimPrefix(*release.Name, configItem.ReleasePrefix+"-")
// 			v, err := semver.NewVersion(versionString)
// 			if err != nil {
// 				log.Infof("Error parsing %s version: %s", versionString, err)
// 			} else if v.Prerelease() == "" {
// 				vs = append(vs, v)
// 			}
// 		}

// 		// stop loop if 10 proper releases are found
// 		if len(vs) >= 10 {
// 			break
// 		}
// 	}

// 	log.Infof("Found %v semver releases with prefix %v for repo %v", len(vs), configItem.ReleasePrefix, *repo.Name)
// 	sort.Sort(sort.Reverse(semver.Collection(vs)))
// 	for _, ver := range vs {
// 		fmt.Printf("\t%s\n", ver.String())
// 	}

// 	return vs
// }

func findLatestVersions(currentVersion *semver.Version, githubVersions []*semver.Version) latestVersions {
	latestMajorVersion := githubVersions[0]
	var latestMinorVersion *semver.Version
	var latestPatchVersion *semver.Version

	for i := 0; i < len(githubVersions); i++ {
		if latestMinorVersion == nil && githubVersions[i].Major() == currentVersion.Major() {
			latestMinorVersion = githubVersions[i]
		}

		if latestPatchVersion == nil && githubVersions[i].Major() == currentVersion.Major() && githubVersions[i].Minor() == currentVersion.Minor() {
			latestPatchVersion = githubVersions[i]
		}

		if latestMinorVersion != nil && latestPatchVersion != nil {
			break
		}
	}

	return latestVersions{
		majorVersion: latestMajorVersion,
		minorVersion: latestMinorVersion,
		patchVersion: latestPatchVersion,
	}
}
