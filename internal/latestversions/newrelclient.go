package latestversions

// func GetAppsLatestVersionVersions() []data.AppVersionInfo{
// allProjects := GetAllProjects()
// }

// func (newreleases.project []newreleases.Project) GetReleases() []newreleases.Release {
// client := newreleases.NewClient(NewReleasesApiKey, nil)
// ctx := context.Background()

// page := 1

// for {
// releases, lastPage, err := client.Releases.ListByProjectID(ctx, project.ID, 1)
// if err != nil {
// log.Fatalln(err)
// }

// if page >= lastPage {
// break
// }
// page++

// return releases
// }
// }

// func parseReleaseToSemver()  {
// vs := []*semver.Version{}
// for _, r := range vs {
// v, err := semver.NewVersion(r.String())
// if err != nil {
// fmt.Printf("Error parsing version: %s\n", err)
// } else if v.Prerelease() == "" {
// vs = append(vs, v)
// }
// }

// sort.Sort(sort.Reverse(semver.Collection(vs)))
// latestMajorVersion := vs[0]
// var latestMinorVersion *semver.Version
// var latestPatchVersion *semver.Version

// for i := 0; i < len(vs); i++ {
// if latestMinorVersion == nil && vs[i].Major() == currentVersion.Major() {
// latestMinorVersion = vs[i]
// }

// if latestPatchVersion == nil && vs[i].Major() == currentVersion.Major() && vs[i].Minor() == currentVersion.Minor() {
// latestPatchVersion = vs[i]
// }
// }

// return data.AppVersionInfo{
// AppName:            repository,
// CurrentVersion:     currentVersion,
// LatestMajorVersion: latestMajorVersion,
// LatestMinorVersion: latestMinorVersion,
// LatestPatchVersion: latestPatchVersion,
// }
// }
