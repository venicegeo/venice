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
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"io"
)

const FIELD_ID = "#"
const FIELD_PROJECT = "Project"
const FIELD_TRACKER = "Tracker"
const FIELD_PARENT = "Parent task"
const FIELD_STATUS = "Status"
const FIELD_PRIORITY = "Priority"
const FIELD_SUBJECT = "Subject"
const FIELD_AUTHOR = "Author"
const FIELD_ASSIGNEE = "Assignee"
const FIELD_UPDATED = "Updated"
const FIELD_CATEGORY = "Category"
const FIELD_TARGET_VERSION = "Target version"
const FIELD_START_DATE = "Start date"
const FIELD_DUE_DATE = "Due date"
const FIELD_ESTIMATED_TIME = "Estimated time"
const FIELD_TOTAL_ESTIMATED_TIME = "Total estimated time"
const FIELD_SPENT_TIME = "Spent time"
const FIELD_TOTAL_SPENT_TIME = "Total spent time"
const FIELD_PERCENT_DONE = "% Done"
const FIELD_CREATED = "Created"
const FIELD_CLOSED = "Closed"
const FIELD_RELATED_ISSUES = "Related issues"
const FIELD_PRIVATE = "Private"

type Issue map[string]string
type Issues map[string]Issue

var fields map[string]string
var data Issues

func assert(id string, condition bool, mssg string) {
	if !condition {
		fmt.Printf("%s: %s -- FAILED\n", id, mssg)
	}
}

func readFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	fields, err := reader.Read()
	if err == io.EOF {
		log.Fatal("no header row found")
	}
	if err != nil {
		log.Fatal(err)
	}
	
	data := make(Issues)
	
	for {
		values, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		issue := make(map[string]string)
		for i, v := range(values) {
			field := fields[i]
			issue[field] = v
		}

		id := issue[FIELD_ID]
		(data)[id] = issue
	}
	
	fmt.Printf("%d columns, %d rows\n", len(fields), len(data))
}

func isEpic(issue Issue) bool {
	return issue[FIELD_CATEGORY] == "Epic"
}

func isStory(issue Issue) bool {
	return issue[FIELD_CATEGORY] == "Story"
}

func isTask(issue Issue) bool {
	return issue[FIELD_CATEGORY] == "Task"
}

func isFuture(issue Issue) bool {
	return issue[FIELD_CATEGORY] == "Future"
}

func parent(issue Issue) Issue {
		id := issue[FIELD_PARENT]
		if id == "" {
			return nil
		}
		return data[id]
}


func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage:  $ redmine-scrub issues.csv")
	}
	
	filename := os.Args[1]
	readFile(filename)
	
	/*for _, v := range(fields) {
		fmt.Printf("%s\n", v)
	}*/
/*	for k := range(*data) {
		fmt.Printf("%s\n", k)
	}*/
	
	//fmt.Print(id,issue,"\n")

	for id, issue := range(data) {
		assert(id, issue[FIELD_CATEGORY] != "", "category is not empty")

		switch issue[FIELD_CATEGORY] {
		case "Epic":
			assert(id, parent(issue) == nil, "epic's parent is nil")
		case "Story":
			assert(id, parent(issue) != nil, "story's parent is not nil")
		case "Task":
			assert(id, parent(issue) == nil, "tasks's parent is not nil")
		case "Future":
			assert(id, parent(issue) == nil, "future's parent is nil")
		default:
			assert(id, false, "category is set")
		}
		
		if isStory(issue) || isTask(issue) {
			assert(id, issue[FIELD_PARENT] != "", "parent is not be empty")
		} else {
			assert(id, issue[FIELD_PARENT] == "", "parent is empty")			
		}
	}
}
