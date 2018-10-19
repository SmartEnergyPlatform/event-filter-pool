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
	"encoding/json"
	"github.com/SmartEnergyPlatform/event-filter-pool/util"
	"net/http"
	"strings"
)

type CamundaEventMsgVariable struct {
	Value interface{} `json:"value,omitempty"`
}

type CamundaEventMsg struct {
	MessageName     string                             `json:"messageName"`
	ProcessVariable map[string]CamundaEventMsgVariable `json:"processVariables"`
	All             bool                               `json:"all"`
	ResultEnabled   bool                               `json:"resultEnabled"`
}

func TriggerEvent(eventid string, msg string) error {
	eventid = strings.Replace(eventid, "#", "_", -1)
	var val interface{}
	err := json.Unmarshal([]byte(msg), &val)
	if err != nil {
		return err
	}
	event := CamundaEventMsg{
		All:         true,
		MessageName: eventid,
		ProcessVariable: map[string]CamundaEventMsgVariable{
			"event": {
				Value: val,
			},
		},
		ResultEnabled: true,
	}
	payload := new(bytes.Buffer)
	err = json.NewEncoder(payload).Encode(event)
	if err != nil {
		return err
	}
	_, err = http.Post(util.Config.CamundaUrl+"/message", "application/json", payload)
	return err
}
