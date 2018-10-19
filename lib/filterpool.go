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
	"log"
	"sync"
)

type FilterCollection struct {
	mux        sync.Mutex
	idIndex    map[string]*Filter            //filterid
	routeIndex map[string]map[string]*Filter //safeConcat(device.service.topic).filterid
	size       int
}

func safeConcat(device string, service string, topic string) string {
	result, _ := json.Marshal([]string{device, service, topic})
	return string(result)
}

var sessionsCollection *FilterCollection
var onceSessionsCollection sync.Once

func FilterPool() *FilterCollection {
	onceSessionsCollection.Do(func() {
		ClearPts()
		sessionsCollection = &FilterCollection{
			idIndex:    map[string]*Filter{},
			routeIndex: map[string]map[string]*Filter{},
			size:       0,
		}
	})
	return sessionsCollection
}

func (this *FilterCollection) Dispatch(topic string, device string, service string, msg string) {
	this.mux.Lock()
	key := safeConcat(device, service, topic)
	filters, ok := this.routeIndex[key]
	this.mux.Unlock()
	if !ok {
		log.Println("WARNING: no matching filter found (", key, ")")
	}
	for _, filter := range filters {
		if ok, jsonValue := filter.Check(msg); ok {
			err := TriggerEvent(filter.Id, jsonValue)
			if err != nil {
				log.Println("ERROR while triggering event: ", err)
			}
		}
	}
}

func (this *FilterCollection) Register(filter *Filter) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.idIndex[filter.Id] = filter
	if _, ok := this.routeIndex[safeConcat(filter.DeviceId, filter.ServiceId, filter.Topic)]; !ok {
		err = RegisterPts(filter.DeviceId, filter.ServiceId, filter.Topic)
		if err != nil {
			return err
		}
		this.routeIndex[safeConcat(filter.DeviceId, filter.ServiceId, filter.Topic)] = map[string]*Filter{}
	}
	this.routeIndex[safeConcat(filter.DeviceId, filter.ServiceId, filter.Topic)][filter.Id] = filter
	this.size++
	return nil
}

func (this *FilterCollection) Deregister(filterId string) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	filter, ok := this.idIndex[filterId]
	if !ok {
		log.Println("WARNING: Deregister() filter to deregister not found ", filterId)
	} else {
		key := safeConcat(filter.DeviceId, filter.ServiceId, filter.Topic)
		if len(this.routeIndex[key]) == 1 {
			err = DeregisterPts(filter.DeviceId, filter.ServiceId, filter.Topic)
			if err != nil {
				return err
			}
			delete(this.idIndex, filterId)
			delete(this.routeIndex[key], filter.Id)
			delete(this.routeIndex, key)
		} else {
			delete(this.idIndex, filterId)
			delete(this.routeIndex[key], filter.Id)
		}
	}
	this.size--
	return nil
}

func (this *FilterCollection) GetSize() int {
	return this.size
}

func (this *FilterCollection) Reset() {
	ClearPts()
	this.mux.Lock()
	this.idIndex = map[string]*Filter{}
	this.routeIndex = map[string]map[string]*Filter{}
	this.size = 0
	this.mux.Unlock()
}
