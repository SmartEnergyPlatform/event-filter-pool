/*
 * Copyright 2019 InfAI (CC SES)
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

package kafka

import (
	"context"
	"sync"
)

type KafkaMultiConsumer struct {
	zkUrl        string
	groupId      string
	ctx          context.Context
	cancel       context.CancelFunc
	consumer     map[string]*KafkaConsumer
	listener     func(topic string, msg []byte) error
	mux          sync.Mutex
	errIn        chan error
	errorhandler func(error)
}

func NewKafkaMultiConsumer(zookeeperUrl string, groupId string, listener func(topic string, msg []byte) error, errorhandler func(error)) (consumer *KafkaMultiConsumer) {
	consumer = &KafkaMultiConsumer{consumer: map[string]*KafkaConsumer{}, listener: listener, zkUrl: zookeeperUrl, groupId: groupId, errorhandler: errorhandler}
	consumer.Init()
	return
}

func (this *KafkaMultiConsumer) Init() {
	this.ctx, this.cancel = context.WithCancel(context.Background())
	this.startErrorHandler()
}

func (this *KafkaMultiConsumer) startErrorHandler() {
	this.errIn = make(chan error)
	go func() {
		done := this.ctx.Done()
		errIn := this.errIn
		for {
			select {
			case <-done:
				return
			case err := <-errIn:
				this.errorhandler(err)
			}
		}
	}()
}

func (this *KafkaMultiConsumer) Listen(topic string) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if consumer, ok := this.consumer[topic]; ok {
		consumer.ListenerCountUp()
	} else {
		consumer, err = NewKafkaConsumer(this.ctx, this.errIn, this.zkUrl, this.groupId, topic, this.listener)
		if err == nil {
			this.consumer[topic] = consumer
		}
	}
	return
}

func (this *KafkaMultiConsumer) Mute(topic string) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if consumer, ok := this.consumer[topic]; ok {
		remove := consumer.ListenerCountDown()
		if remove {
			consumer.Stop()
			delete(this.consumer, topic)
		}
	}
}

func (this *KafkaMultiConsumer) Reset() {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.cancel()
	this.ctx, this.cancel = context.WithCancel(context.Background())
	this.startErrorHandler()
	for _, consumer := range this.consumer {
		consumer.Start(this.ctx)
	}
}

func (this *KafkaMultiConsumer) Close() {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.cancel()
}
