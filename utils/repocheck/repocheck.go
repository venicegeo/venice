/*
Copyright 2016, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This app reads in a CSV export of a Redmine project and "validates" all
// the issues based on various rules such as:
//    a Story must have a parent, and that parent must be an Epic
//    a Task in the current milestone must have an estimate, and that
//      estimate must be <= 16 hours
// This app is still underdevelopment -- will be adding rules as we need them.
//
// To use:
//   - in Redmine, click on "View all issues" (right-hand column, top)
//   - then click on "Also available in... CSV" (main panel, bottom-right)
//   - run the app:  $ go run redmine-scrub.go downloadedfile.csv
//

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

const DEBUG = false
const VENICE = "venicegeo"

var NOW = time.Now()
var WHO = "mpgerlek"
var EMAIL = "mpg@flaxen.com"

var reposWhiteList = []string{
	"bf-handle",
	"bf_TidePrediction",
	"geojson-geos-go",
	"geojson-go",
	"pz-access",
	"pz-docs",
	"pz-gateway",
	"pz-gocommon",
	"pz-idam",
	"pz-ingest",
	"pz-jobcommon",
	"pz-jobmanager",
	"pz-logger",
	"pz-sak",
	"pz-search",
	"pz-servicecontroller",
	"pz-swagger",
	"pz-uuidgen",
	"pz-workflow",
	"pzsvc-exec",
	"pzsvc-image-catalog",
	"pzsvc-lib",
	"pzsvc-ossim",
}

var extWhiteList = []string{
	".go",
	".lock",
	".md",
	".sh",
	".txt",
	".yaml",
	".yml",
}

func contains(array []string, item string) bool {
	for _, i := range array {
		if i == item {
			return true
		}
	}
	return false
}

func main() {
	//log.Printf("%#v", os.Args)

	if len(os.Args) == 2 && os.Args[1] == "-update" {
		DoUpdate()
	} else if len(os.Args) == 3 && os.Args[1] == "-check" {
		err := DoCheck(os.Args[2])
		if err != nil {
			log.Fatalf("ERROR: %s", err.Error())
		}
	} else {
		s := `usage:  $ repocheck -update
        $ repocheck -check <reponame>
`
		fmt.Printf(s)
		os.Exit(1)
	}

}

func DoCheck(repoName string) error {
	f, err := os.Open(repoName)
	if err != nil {
		return err
	}

	files, err := f.Readdirnames(0)
	if err != nil {
		return err
	}

	f.Close()

	if !contains(files, "README") &&
		!contains(files, "README.txt") &&
		!contains(files, "README.md") {
		fmt.Printf("%s: no README", repoName)
	}

	if !contains(files, "LICENSE") &&
		!contains(files, "LICENSE.txt") {
		fmt.Printf("%s: no LICENSE", repoName)
	}

	err = inspectDirectory(repoName)
	if err != nil {
		return err
	}

	return nil
}

func isDotFile(name string) bool {
	_, ff := filepath.Split(name)
	return ff[0] == '.'
}

func inspectDirectory(dirName string) error {
	//fmt.Printf("...d %s\n", dirName)

	if isDotFile(dirName) {
		return nil
	}

	f, err := os.Open(dirName)
	if err != nil {
		return err
	}

	files, err := f.Readdirnames(0)
	if err != nil {
		return err
	}

	f.Close()

	for _, file := range files {
		inspectFile(dirName + "/" + file)
	}

	return nil
}

func inspectFile(fileName string) error {
	//fmt.Printf("...f %s\n", fileName)

	if isDotFile(fileName) {
		return nil
	}

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}
	f.Close()

	if fileInfo.IsDir() {
		return inspectDirectory(fileName)
	} else {
		ext := filepath.Ext(fileName)
		if ext != "" {

			if !contains(extWhiteList, ext) {
				fmt.Printf("%s: unknown suffix\n", fileName)
			}
		}
	}

	return nil
}

func DoUpdate() {
	apiKey, err := getApiKey()
	if err != nil {
		log.Fatalf("Failed to get API key: %s", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiKey},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	//log.Printf("client: %#v", client)

	// list all repositories for the authenticated user
	repos, err := getRepoNames(client)

	for _, repo := range repos {
		if contains(reposWhiteList, *repo.FullName) {
			exec.Command("ls")
			// if present, update
			// else download
		} else if contains(reposBlackList, *repo.FullName) {
			// if present, error out
		} else {
			// new repo name, error out
		}
		fmt.Printf("    \"%s\",\n", *repo.FullName)
	}
}

func getApiKey() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("getApiKey: $HOME not found")
	}

	key, err := ioutil.ReadFile(home + "/.git-token")
	if err != nil {
		return "", fmt.Errorf("getApiKey: %s", err)
	}

	s := strings.TrimSpace(string(key))
	//log.Printf("API Key: %s", s)

	return s, nil
}

func getRepoNames(client *github.Client) ([]*github.Repository, error) {
	opts0 := github.ListOptions{
		Page:    0,
		PerPage: 512,
	}
	opts := &github.RepositoryListByOrgOptions{
		ListOptions: opts0,
		Type:        "all",
	}

	repos, _, err := client.Repositories.ListByOrg(VENICE, opts)
	if err != nil {
		return nil, fmt.Errorf("getRepoNames: %s", err)
	}

	//log.Printf("%#v", repos)

	return repos, nil
}
