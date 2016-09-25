package com.lovi.quebic.sockm.messenger;

import com.lovi.quebic.sockm.message.messenger.Message;

public interface Publisher {

	<I,O> void addSubscriber(String subscriberId, Subscriber<I,O> subscriber);

	<I> void publish(String address, I message);

	/**
	 * If found a subscriber for the address return true. Otherwise return false
	 * @param address
	 * @param message
	 * @return
	 */
	<I, O> boolean send(String address, Message<I, O> message);

	<I, O> boolean sendBySubscriberId(String subscriberId, Message<I, O> message);

}
