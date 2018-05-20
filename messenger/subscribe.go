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
)

//Subscribe subscribe
func (messenger *Messenger) Subscribe(eventID string, requestHandler func(baseEvent BaseEvent), consumerName string) error {

	if eventID == "" {
		return fmt.Errorf("eventID should not be empty")
	}

	if requestHandler == nil {
		return fmt.Errorf("requestHandler should not be nil")
	}

	//eventID become the routingKey
	routingKey := eventID
	queueName := messenger.getQueueName(routingKey)
	consumerID, _ := createConsumerID(consumerName)

	queue, err := messenger.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		true,      // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare %s queue, error : %v", queueName, err)
	}

	err = messenger.channel.QueueBind(
		queue.Name, // queue name
		routingKey, // routing key
		Exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind %s queue, error : %v", queueName, err)
	}

	msgs, err := messenger.channel.Consume(
		queue.Name, // queue
		consumerID, // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("failed to consume %s queue, error : %v", queueName, err)
	}

	go func() {

		for msg := range msgs {

			baseEvent := BaseEvent{eventPayload: msg.Body, headers: msg.Headers}

			//TODO remove later
			//log.Printf("received event : %s", baseEvent.GetEventID())

			requestHandler(baseEvent)

		}

	}()

	//TODO remove log
	log.Printf("subscribed for event : %s", routingKey)

	return nil

}
