/*
Copyright 2018 Tharanga Nilupul Thennakoon

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
package components

import (
	"fmt"
	"quebic-faas/config"
	"time"

	"github.com/streadway/amqp"
)

//WaitForEventbusConnect waitForEventbusConnect
func WaitForEventbusConnect(eventBusConfig config.EventBusConfig, wait time.Duration) error {

	waitForResponse := make(chan bool)

	go func() {

		connectionStr := eventBusConnectionStr(eventBusConfig)

		for {

			_, err := amqp.Dial(connectionStr)
			if err == nil {
				waitForResponse <- true
				break
			}

			time.Sleep(time.Second * 1)

		}

	}()

	select {
	case <-waitForResponse:
		break
	case <-time.After(wait):
		return fmt.Errorf("unable to connect eventbus")
	}

	return nil

}

//WaitForEventbusStop waitForEventbusStop
func WaitForEventbusStop(eventBusConfig config.EventBusConfig, wait time.Duration) error {

	waitForResponse := make(chan bool)

	go func() {

		connectionStr := eventBusConnectionStr(eventBusConfig)

		for {

			_, err := amqp.Dial(connectionStr)
			if err != nil {
				waitForResponse <- true
				break
			}

			time.Sleep(time.Second * 1)

		}

	}()

	select {
	case <-waitForResponse:
		break
	case <-time.After(wait):
		return fmt.Errorf("previous eventbus still running on. please try again")
	}

	return nil

}

func eventBusConnectionStr(eventBusConfig config.EventBusConfig) string {
	var connectionStr string
	if eventBusConfig.AMQPPort == 0 {
		connectionStr = fmt.Sprintf("amqp://%s:%s@%s/",
			eventBusConfig.ManagementUserName,
			eventBusConfig.ManagementPassword,
			eventBusConfig.AMQPHost)
	} else {
		connectionStr = fmt.Sprintf("amqp://%s:%s@%s:%d/",
			eventBusConfig.ManagementUserName,
			eventBusConfig.ManagementPassword,
			eventBusConfig.AMQPHost,
			eventBusConfig.AMQPPort)
	}
	return connectionStr
}
