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
	"github.com/segmentio/kafka-go"
	"io"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

func NewKafkaConsumer(ctx context.Context, errOut chan<- error, zk string, groupid string, topic string, listener func(topic string, msg []byte) error) (consumer *KafkaConsumer, err error) {
	consumer = &KafkaConsumer{count: 0, groupId: groupid, zkUrl: zk, err: errOut, topic: topic, listener: listener}
	err = consumer.Start(ctx)
	return
}

type KafkaConsumer struct {
	count    int
	zkUrl    string
	groupId  string
	topic    string
	ctx      context.Context
	cancel   context.CancelFunc
	err      chan<- error
	listener func(topic string, msg []byte) error
	mux      sync.Mutex
}

func (this *KafkaConsumer) ListenerCountUp() {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.count = this.count + 1
}

func (this *KafkaConsumer) ListenerCountDown() (remove bool) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.count = this.count - 1
	return this.count < 1
}

func (this *KafkaConsumer) Stop() {
	this.cancel()
}

func (this *KafkaConsumer) Start(ctx context.Context) error {
	log.Println("DEBUG: consume topic: \"" + this.topic + "\"")
	this.ctx, this.cancel = context.WithCancel(ctx)
	broker, err := GetBroker(this.zkUrl)
	if err != nil {
		log.Println("ERROR: unable to get broker list", err)
		return err
	}
	err = InitTopic(broker, this.topic)
	if err != nil {
		log.Println("ERROR: unable to create topic", err)
		return err
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     broker,
		GroupID:     this.groupId,
		Topic:       this.topic,
		MaxWait:     1 * time.Second,
		Logger:      log.New(ioutil.Discard, "", 0),
		ErrorLogger: log.New(ioutil.Discard, "", 0),
	})
	go func() {
		for {
			select {
			case <-this.ctx.Done():
				log.Println("close kafka reader ", this.topic)
				return
			default:
				m, err := r.FetchMessage(this.ctx)
				if err == io.EOF || err == context.Canceled {
					log.Println("close consumer for topic ", this.topic)
					return
				}
				if err != nil {
					log.Println("ERROR: while consuming topic ", this.topic, err)
					this.err <- err
					return
				}
				err = this.listener(m.Topic, m.Value)
				if err != nil {
					log.Println("ERROR: unable to handle message (no commit)", err)
				} else {
					//log.Println("DEBUG: commit for ", m.Topic)
					err = r.CommitMessages(this.ctx, m)
				}
			}
		}
	}()
	return err
}
