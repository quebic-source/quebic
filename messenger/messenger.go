//    Copyright 2018 Tharanga Nilupul Thennakoon
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package messenger

import (
	"fmt"
	"log"
	"net/http"
	"quebic-faas/common"
	"quebic-faas/config"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

const Exchange = "quebic-faas-eventbus"

const emptyError = ""
const defaultStatuscode = 200

//Messenger interface
type Messenger struct {
	AppID          string
	EventBusConfig config.EventBusConfig
	connection     *amqp.Connection
	channel        *amqp.Channel
}

//Init create connection
func (messenger *Messenger) Init() error {

	if messenger.connection == nil {

		config := messenger.EventBusConfig

		var connectionStr string
		if config.AMQPPort == 0 {
			connectionStr = fmt.Sprintf("amqp://%s:%s@%s/",
				config.ManagementUserName,
				config.ManagementPassword,
				config.AMQPHost)
		} else {
			connectionStr = fmt.Sprintf("amqp://%s:%s@%s:%d/",
				config.ManagementUserName,
				config.ManagementPassword,
				config.AMQPHost,
				config.AMQPPort)
		}

		connection, err := amqp.Dial(connectionStr)
		if err != nil {
			return fmt.Errorf("failed to connect to eventbus, error : %v", err)
		}
		messenger.connection = connection

		log.Printf("quebic-faas-eventbus : connected")

	}

	if messenger.channel == nil {

		channel, err := messenger.connection.Channel()
		if err != nil {
			return fmt.Errorf("failed to open eventbus channel, error : %v", err)
		}
		messenger.channel = channel

		log.Printf("quebic-faas-eventbus : create channel")

	}

	err := messenger.channel.ExchangeDeclare(
		Exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return fmt.Errorf("failed to declare %s exchange, error : %v", Exchange, err)
	}

	log.Printf("quebic-faas-eventbus : create exchange")

	return nil

}

//WaitInit create connection
func (messenger *Messenger) WaitInit(wait time.Duration) error {

	waitForResponse := make(chan bool)

	go func() {

		log.Printf("waiting for connecting to eventbus...")

		for {

			err := messenger.Init()
			if err == nil {

				err = messenger.pingToRabbitmqMgr()
				if err == nil {
					waitForResponse <- true
					break
				}

			}

			time.Sleep(time.Second * 2)

		}
	}()

	select {
	case <-waitForResponse:
		return nil
	case <-time.After(wait):
		return fmt.Errorf("eventbus init request failed. reason : timeout")
	}

}

//Close close connection
func (messenger *Messenger) Close() {

	if messenger.channel != nil {
		err := messenger.channel.Close()
		if err != nil {
			log.Printf("failed to close eventbus channel, error : %v", err)
		}
	}

	if messenger.connection != nil {
		err := messenger.connection.Close()
		if err != nil {
			log.Printf("failed to close eventbus connection, error : %v", err)
		}
	}

}

//ReleseQueue release queue
func (messenger *Messenger) ReleseQueue(routingKey string) {

	err := messenger.channel.Cancel(routingKey, false)
	if err != nil {
		log.Printf("failed to cancel consumer %s, error : %v", routingKey, err)
	}

	_, err = messenger.channel.QueueDelete(
		messenger.getQueueName(routingKey),
		false, // if unused
		false, // if unused
		false, // no wait
	)
	if err != nil {
		log.Printf("failed to releseQueue %s, error : %v", routingKey, err)
	}

}

func (messenger *Messenger) getQueueName(routingKey string) string {
	//when running on function, functionID become the AppID
	return fmt.Sprintf("%s.%s", messenger.AppID, routingKey)
}

func createRequestID(baseEvent *BaseEvent) error {

	requestIDUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("unable to assign requestID for baseEvent %v", err)
	}
	requestID := requestIDUUID.String()

	baseEvent.setRequestID(requestID)

	return nil
}

func createExecutionID(baseEvent *BaseEvent) error {

	requestIDUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("unable to assign execution-ack id for baseEvent %v", err)
	}
	requestID := requestIDUUID.String()

	baseEvent.setRequestID(createEventID(common.EventPrefixExecutionACK, requestID))

	return nil
}

func createEventID(eventType string, eventIDSuffix string) string {
	return eventType + common.EventJOIN + eventIDSuffix
}

func createConsumerID(consumerName string) (string, error) {

	uuid, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("unable to assign requestID for baseEvent %v", err)
	}

	return consumerName + common.EventJOIN + uuid.String(), nil
}

func (messenger *Messenger) pingToRabbitmqMgr() error {

	constr := fmt.Sprintf(
		"http://%s:%d",
		messenger.EventBusConfig.ManagementHost,
		messenger.EventBusConfig.ManagementPort,
	)

	_, err := http.Get(constr)
	if err != nil {
		return err
	}

	return nil

}
