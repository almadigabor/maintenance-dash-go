package data

import "github.com/Masterminds/semver/v3"

type AppVersionInfo struct {
	AppName            string
	CurrentVersion     *semver.Version
	LatestMajorVersion *semver.Version
	LatestMinorVersion *semver.Version
	LatestPatchVersion *semver.Version
}
