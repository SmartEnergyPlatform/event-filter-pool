/*
 * Copyright 2018 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lib

import (
	"encoding/json"
	"github.com/SmartEnergyPlatform/event-filter-pool/util"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func StartManagerConnection() {
	log.Println("StartManagerConnection")
	InitAssignments()
	go func() {
		for {
			wait := NextManagerCommand()
			if wait {
				duration := time.Duration(util.Config.ManagerCommandGetTimeout) * time.Millisecond
				time.Sleep(duration)
			}
		}
	}()
}

func InitAssignments() {
	assignments, err := getAllAssignments()
	if err != nil {
		log.Fatalln("ERROR: InitAssignments() ", err)
	}
	handleCommand("add", assignments)
}

func NextManagerCommand() (wait bool) {
	command, err := getCommands()
	if err != nil {
		log.Println("ERROR: on getCommands", err)
		return true
	}
	handleCommand(command.Command, command.Assignment)
	return command.Command == "none"
}

func getAllAssignments() (assignments []FilterDeployment, err error) {
	resp, err := http.Get(util.Config.FilterManagerUrl + "/pool/assignments/" + util.Config.PoolId)

	if err != nil {
		log.Println("error on getAllAssignments", err)
		return assignments, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&assignments)
	return
}

func getCommands() (command FilterPoolCommand, err error) {
	size := strconv.Itoa(FilterPool().GetSize())
	resp, err := http.Get(util.Config.FilterManagerUrl + "/pool/command/" + util.Config.PoolId + "/" + size)

	if err != nil {
		log.Println("error on getCommands", err)
		return command, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&command)
	return
}

func handleCommand(command string, assignments []FilterDeployment) {
	switch command {
	case "reset":
		log.Println("reset: ", assignments)
		FilterPool().Reset()
		handleCommand("add", assignments)
	case "add":
		log.Println("add: ", assignments)
		for _, assignment := range assignments {
			filter, err := NewFilter(assignment)
			if err != nil {
				log.Println("ERROR: ", err)
				ReportErrorToManager(assignment.FilterId, err)
				return
			}
			err = FilterPool().Register(&filter)
			if err != nil {
				log.Println("ERROR: ", err)
				ReportErrorToManager(assignment.FilterId, err)
				return
			}
		}
	case "remove":
		log.Println("remove: ", assignments)
		for _, assignment := range assignments {
			err := FilterPool().Deregister(assignment.FilterId)
			if err != nil {
				log.Println("ERROR: ", err)
				ReportErrorToManager(assignment.FilterId, err)
				return
			}
		}
	case "none":

	default:
		log.Println("WARNING: unknown command ", command)

	}
}

func ReportErrorToManager(filterid string, err error) {
	resp, err := http.Get(util.Config.FilterManagerUrl + "/pool/error/" + util.Config.PoolId + "/" + url.QueryEscape(filterid) + "/" + url.QueryEscape(err.Error()))
	if err != nil {
		log.Println("ERROR: on report error to manager", err)
		return
	}
	defer resp.Body.Close()
}
