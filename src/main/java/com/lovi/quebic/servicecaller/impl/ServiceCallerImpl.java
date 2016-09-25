package com.lovi.quebic.servicecaller.impl;

import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.servicecaller.MessageBody;
import com.lovi.quebic.servicecaller.ServiceCaller;
import com.lovi.quebic.sockm.launcher.SockMLauncher;
import com.lovi.quebic.sockm.messenger.Messenger;

public class ServiceCallerImpl<I,O> implements ServiceCaller<I,O>{

	private Messenger messenger;
	
	private String address;
	
	private Class<?> inputClassType;
	
	private Class<?> outputClassType;
	
	private Handler<O> resultHandler;
	
	private Handler<Throwable> failureHandler;
	
	public ServiceCallerImpl(String address){
		this(address, MessageBody.class, Object.class);
	}
	
	public ServiceCallerImpl(String address, Class<?> inputClassType, Class<?> outputClassType){
		this.messenger = SockMLauncher.create().getMessenger();
		this.address = address;
		this.inputClassType = inputClassType;
		this.outputClassType = outputClassType;
	}
	
	@Override
	public ServiceCaller<I,O> result(Handler<O> resultHandler){
		this.resultHandler = resultHandler;
		return this;
	}
	
	@Override
	public ServiceCaller<I,O> failure(Handler<Throwable> failureHandler){
		this.failureHandler = failureHandler;
		return this;
	}
	
	@Override
	public void call(Object... objects){
		if(inputClassType.getName().equals(MessageBody.class.getName())){
			if(resultHandler != null)
				messenger.send(address, new MessageBody(objects), inputClassType, outputClassType, resultHandler, failureHandler);
			else
				messenger.publish(address, new MessageBody(objects));
		}else{
			if(resultHandler != null)
				messenger.send(address, objects[0], inputClassType, outputClassType, resultHandler, failureHandler);
			else
				messenger.publish(address, objects[0]);
		}
	}
	
}
