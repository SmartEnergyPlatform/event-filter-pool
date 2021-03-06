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
	"bytes"
	"errors"
	"github.com/SmartEnergyPlatform/event-filter-pool/util"
	"log"
	"net/http"
	"net/url"
)

func RegisterPts(device string, service string, topic string) (err error) {
	log.Println("register route: ", device, service, topic, util.Config.PoolId)
	resp, err := http.Post(util.Config.PtsUrl+"/add/route/"+url.QueryEscape(topic)+"/"+url.QueryEscape(device)+"/"+url.QueryEscape(service)+"/"+util.Config.PoolId, "", nil)

	if err != nil {
		log.Println("error on RegisterPrefixTopic", err)
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()

	if s != "ok" {
		return errors.New("unexpected pts response: " + s)
	}
	return
}

func DeregisterPts(device string, service string, topic string) (err error) {
	log.Println("deregister route: ", device, service, topic, util.Config.PoolId)
	req, err := http.NewRequest("DELETE", util.Config.PtsUrl+"/remove/route/"+url.QueryEscape(topic)+"/"+url.QueryEscape(device)+"/"+url.QueryEscape(service)+"/"+util.Config.PoolId, nil)
	if err != nil {
		log.Println("error while building DeregisterPts() request: ", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("error on DeregisterPts()", err)
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()

	if s != "ok" {
		return errors.New("unexpected pts response: " + s)
	}
	return
}

func ClearPts() (err error) {
	req, err := http.NewRequest("DELETE", util.Config.PtsUrl+"/remove/target/"+util.Config.PoolId, nil)
	if err != nil {
		log.Println("error while building DeregisterPts() request: ", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("error on ClearPts()", err)
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()

	if s != "ok" {
		return errors.New("unexpected pts response: " + s)
	}
	return
}
