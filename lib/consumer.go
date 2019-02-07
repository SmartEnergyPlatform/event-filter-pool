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
	"github.com/SENERGY-Platform/iot-broker-client"
	"github.com/SmartEnergyPlatform/event-filter-pool/util"
	"log"

	"encoding/json"

)

func InitConsumer() (consumer *iot_broker_client.Consumer,err error) {
	consumer, err = iot_broker_client.NewConsumer(util.Config.AmqpUrl, util.Config.PoolId, util.Config.FilterTopic, false, func(msg []byte) error {
		go HandleMessage(util.Config.FilterTopic, string(msg))
		return nil
	})
	return
}

type Envelope struct {
	DeviceId  string      `json:"device_id,omitempty"`
	ServiceId string      `json:"service_id,omitempty"`
	Value     interface{} `json:"value"`
}

func HandleMessage(topic string, msg string) {
	log.Println("consume kafka msg: ", msg)
	envelope := Envelope{}
	err := json.Unmarshal([]byte(msg), &envelope)
	if err != nil {
		log.Println("ERROR: ", err)
		return
	}
	FilterPool().Dispatch(topic, envelope.DeviceId, envelope.ServiceId, msg)
}

type PrefixMessage struct {
	DeviceId  string      `json:"device_id,omitempty"`
	ServiceId string      `json:"service_id,omitempty"`
	Value     interface{} `json:"value"`
}
