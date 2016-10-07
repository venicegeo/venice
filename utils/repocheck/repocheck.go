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
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const DEBUG = false
const VENICE = "venicegeo"

var NOW = time.Now()
var WHO = "mpgerlek"
var EMAIL = "mpg@flaxen.com"

const giturl = "http://github.com/venicegeo"

var reposWhiteList = []string{
	"bf-handle",
	"bf_TidePrediction",
	"bf-ui",
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
	"pz-metrics",
	"pz-sak",
	"pz-search-metadata-ingest",
	"pz-search-query",
	"pz-servicecontroller",
	//"pz-swagger",		// skip this, not really ours
	"pz-uuidgen",
	"pz-workflow",
	"pzsvc-exec",
	"pzsvc-hello",
	"pzsvc-image-catalog",
	"pzsvc-lib",
	"pzsvc-ossim",
	"pzsvc-lib",
	"pzsvc-ossim",
	"pzsvc-preview-generator",
}

var extCheckList = []string{
	".go",
	".java",
	".js",
	".py",
	".ts",
	".tsx",
}

var extIgnoreList = []string{
	".backup",
	".bak",
	".bat",
	".cmd",
	".conf",
	".config",
	".css",
	".docx",
	".eot",
	".geojson",
	".gif",
	".gz",
	".go",
	".handlebars",
	".html",
	".ico",
	".iml",
	".java",
	".jpg",
	".js",
	".json",
	".laz",
	".less",
	".lock",
	".map",
	".md",
	".off",
	".offf",
	".opts",
	".pdf",
	".per",
	".pkl",
	".png",
	".postman_collection",
	".pptx",
	".properties",
	".py",
	".README",
	".sh",
	".st",
	".svg",
	".tif",
	".ts",
	".tsx",
	".ttf",
	".txt",
	".wkt",
	".woff",
	".woff2",
	".xml",
	".yaml",
	".yml",
	".zip",
}

var specialIgnoreList = []string{
	"bf_TidePrediction/test/__init__.py",
	"bf-ui/node_modules/",
	"bf-ui/src/openlayers.d.ts",
	"pz-sak/public/js/lib/",
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
		err := DoUpdate()
		if err != nil {
			fmt.Printf("aborting: %s\n", err.Error())
		}
	} else if len(os.Args) == 2 && os.Args[1] == "-check" {
		for _, repo := range reposWhiteList {
			err := DoCheck(repo)
			if err != nil {
				fmt.Printf("aborting: %s\n", err.Error())
			}
		}
	} else if len(os.Args) == 3 && os.Args[1] == "-check" {
		err := DoCheck(os.Args[2])
		if err != nil {
			fmt.Printf("aborting: %s\n", err.Error())
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
		fmt.Printf("%s: no README\n", repoName)
	}

	if !contains(files, "LICENSE") &&
		!contains(files, "LICENSE.txt") &&
		!contains(files, "LICENSE.md") {
		fmt.Printf("%s: no LICENSE\n", repoName)
	}

	if !contains(files, ".about.yml") {
		fmt.Printf("%s: no .about.yml\n", repoName)
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

	if strings.HasSuffix(dirName, "/vendor") {
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

func hasCopyright(f *os.File) (bool, error) {
	buf := make([]byte, 4096)

	r := bufio.NewReader(f)
	_, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	s := string(buf)

	ok := strings.Contains(s, "Copyright 2016, RadiantBlue Technologies, Inc.")
	if !ok {
		return false, nil
	}
	ok = strings.Contains(s, "Apache License, Version 2.0")
	if !ok {
		return false, nil
	}
	ok = strings.Contains(s, "WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND")
	if !ok {
		return false, nil
	}

	return true, err
}

func inspectFile(fileName string) error {
	//fmt.Printf("...f %s\n", fileName)

	if isDotFile(fileName) || fileName[len(fileName)-1] == '~' {
		return nil
	}

	for _, ignorable := range specialIgnoreList {
		if strings.Contains(fileName, ignorable) {
			return nil
		}
	}

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}
	defer f.Close()

	if fileInfo.IsDir() {
		return inspectDirectory(fileName)
	}

	ext := filepath.Ext(fileName)
	if ext == "" {
		return nil
	}

	if !contains(extIgnoreList, ext) {
		fmt.Printf("%s: unknown suffix '%s'\n", fileName, ext)
		return nil
	}

	if !contains(extCheckList, ext) {
		return nil
	}

	ok, err := hasCopyright(f)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Printf("%s: no copyright\n", fileName)
	}

	return nil
}

func fileExists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		// other error
		return false, err
	}
	return true, nil
}

func DoUpdate() error {

	for _, repo := range reposWhiteList {

		exists, err := fileExists(repo)
		if err != nil {
			return err
		}

		if exists {
			fmt.Printf("...syncing %s\n", repo)
			out, err := exec.Command("git", "-C", repo, "pull").Output()
			fmt.Printf("OUT: %s\n", out)
			if err != nil {
				fmt.Printf("ERR: %s\n", err.Error())
				return err
			}
		} else {
			fmt.Printf("...cloning %s\n", repo)
			out, err := exec.Command("pwd").Output()
			fmt.Printf("OUT: %s\n", out)
			if err != nil {
				fmt.Printf("ERR: %s\n", err.Error())
				return err
			}
			out, err = exec.Command("git", "clone", giturl+"/"+repo).Output()
			fmt.Printf("OUT: %s\n", out)
			if err != nil {
				fmt.Printf("ERR: %s\n", err.Error())
				return err
			}
		}
	}

	return nil
}
