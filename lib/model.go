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

type Rule struct {
	Path     string `json:"path" bson:"path,omitempty"`         // github.com/NodePrime/jsonpath
	Scope    string `json:"scope" bson:"scope,omitempty"`       // 'any' || 'all' || 'none' || 'max <number>' || 'min <number>' || '<number>'
	Operator string `json:"operator" bson:"operator,omitempty"` // '==' || '!=' || '>=' || '<=' || '<' || '>' || 'regex'
	Value    string `json:"value" bson:"value,omitempty"`
}

type FilterDesc struct {
	Scope     string `json:"scope" bson:"scope,omitempty"` // 'any' || 'all' || 'none' || 'max <number>' || 'min <number>' || '<number>'
	Rules     []Rule `json:"rules" bson:"rules,omitempty"`
	Topic     string `json:"topic" bson:"topic,omitempty"`
	DeviceId  string `json:"device_id,omitempty" bson:"device_id,omitempty"`
	ServiceId string `json:"service_id,omitempty" bson:"device_id,omitempty"`
}

type FilterDeployment struct {
	FilterId   string     `json:"filter_id" bson:"filter_id,omitempty"`
	ProcessId  string     `json:"process_id" bson:"process_id,omitempty"`
	FilterPool string     `json:"filter_pool" bson:"filter_pool,omitempty"`
	State      string     `json:"state" bson:"state,omitempty"`
	Filter     FilterDesc `json:"filter" bson:"filter,omitempty"`
}

type FilterPoolCommand struct {
	Command    string
	Assignment []FilterDeployment
}
