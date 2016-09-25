package com.lovi.quebic.sockm.messenger;

import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.sockm.launcher.SockMLauncher;
import com.lovi.quebic.sockm.message.messenger.Message;
import com.lovi.quebic.sockm.messenger.impl.MessengerImpl;


public interface Messenger {

	static Messenger create(SockMLauncher launcher){
		return MessengerImpl.create(launcher);
	}

	String getLauncherId();

	<I, O> Subscriber<I, O> subscribe(String address, Handler<Message<I, O>> handler);
	
	<I, O> Subscriber<I, O> subscribe(String address, Class<?> inputClassType,
			Class<?> outputClassType, Handler<Message<I, O>> handler);

	<I> void publish(String address, I message);

	<I> void publishLocal(String address, I message);

	<I, O> void send(String address, I message, Handler<O> resultHandler, Handler<Throwable> failureHandler);

	<I, O> void send(String address, I message, Class<?> inputClassType,
			Class<?> outputClassType, Handler<O> resultHandler,
			Handler<Throwable> failureHandler);
	
	<I, O> boolean sendLocal(String subscriberId, Message<I, O> msg);

}
