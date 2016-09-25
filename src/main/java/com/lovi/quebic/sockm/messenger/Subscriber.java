package com.lovi.quebic.sockm.messenger;

import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.sockm.message.messenger.Message;

public interface Subscriber<I, O>{

	Subscriber<I, O> subscribe(String address, Handler<Message<I, O>> handler);
	
	Subscriber<I, O> subscribe(String address, Class<?> inputClassType,
			Class<?> outputClassType, Handler<Message<I, O>> handler);

	Handler<O> getResultHandler();

	void setResultHandler(Handler<O> resultHandler);

	Handler<O> getFailureHandler();

	void setFailureHandler(Handler<O> failureHandler);

	String getAddress();

	Handler<Message<I, O>> getHandler();

	String getId();

	void run(Message<I, O> msg);

	void run(Object message);

}
