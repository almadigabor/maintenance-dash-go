package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ReleaseConfig struct {
	Apps []ReleaseConfigItem `yaml:"apps"`
}

type ReleaseConfigItem struct {
	OrgName       string `yaml:"orgName"`
	RepoName      string `yaml:"repoName"`
	ReleasePrefix string `yaml:"releasePrefix,omitempty"`
}

func readConf(filename string) (*ReleaseConfig, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &ReleaseConfig{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}
