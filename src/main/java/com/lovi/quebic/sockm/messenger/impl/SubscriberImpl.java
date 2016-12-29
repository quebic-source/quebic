package com.lovi.quebic.sockm.messenger.impl;

import java.util.UUID;
import org.codehaus.jackson.map.JsonMappingException;
import org.codehaus.jackson.map.ObjectMapper;

import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.sockm.log.error.ErrorMessage;
import com.lovi.quebic.sockm.message.messenger.Message;
import com.lovi.quebic.sockm.messenger.Publisher;
import com.lovi.quebic.sockm.messenger.Subscriber;

public class SubscriberImpl<I,O> implements Subscriber<I,O>{
	
	private final String id;
	private Publisher publisher;
	private String address;
	private Class<?> inputClassType;
	private Class<?> outputClassType;
	private Handler<Message<I,O>> handler;
	private Handler<O> resultHandler;
	private Handler<O> failureHandler;
	
	public SubscriberImpl(Publisher publisher) {
		this.id = genarateId();
		this.publisher = publisher;
	}

	@Override
	public Subscriber<I,O> subscribe(String address, Handler<Message<I,O>> handler){
		this.address = address;
		this.handler = handler;
		this.publisher.addSubscriber(id, this);
		return this;
	}
	
	@Override
	public Subscriber<I,O> subscribe(String address, Class<?> inputClassType, Class<?> outputClassType, Handler<Message<I,O>> handler){
		this.address = address;
		this.handler = handler;
		this.inputClassType = inputClassType;
		this.outputClassType = outputClassType;
		this.publisher.addSubscriber(id, this);
		return this;
	}
	
	@SuppressWarnings("unchecked")
	@Override
	public void run(Object message){
		try{
			if(handler != null){
				Message<I,O> msg = new Message<>((I)message);
				
				if(inputClassType != null){
					
					ObjectMapper objectMapper = new ObjectMapper();
					I inputMessage = msg.getMessage();
					
					if(!inputMessage.getClass().getName().equals(inputClassType.getName())){
						String jsonStr = objectMapper.writeValueAsString(inputMessage);
						inputMessage = (I) objectMapper.readValue(jsonStr, inputClassType);
						msg.setMessage(inputMessage);
					}
					
				}
				
				handler.handle(msg);
			}
		}catch(Exception e){
		}
		
		
	}
	
	@Override
	public void run(Message<I,O> msg){
		try{
			if(handler != null){
				
				if(inputClassType != null){
					
					ObjectMapper objectMapper = new ObjectMapper();
					I inputMessage = msg.getMessage();
					
					if(!inputMessage.getClass().getName().equals(inputClassType.getName())){
						String jsonStr = objectMapper.writeValueAsString(inputMessage);
						inputMessage = (I) objectMapper.readValue(jsonStr, inputClassType);
						msg.setMessage(inputMessage);
					}
					
				}
				
				handler.handle(msg);
				
			}
		}catch(JsonMappingException e){
			msg.replyFailure(new Exception(ErrorMessage.UNABLE_TO_PARSE_POJO_DEFAULT_CONSTRUCTOR_NOT_FOUND + inputClassType.getName()));
		}catch(Exception e){
			msg.replyFailure(e);
		}
		
	}

	@Override
	public Handler<O> getResultHandler() {
		return resultHandler;
	}

	@Override
	public void setResultHandler(Handler<O> resultHandler) {
		this.resultHandler = resultHandler;
	}

	@Override
	public Handler<O> getFailureHandler() {
		return failureHandler;
	}

	@Override
	public void setFailureHandler(Handler<O> failureHandler) {
		this.failureHandler = failureHandler;
	}

	@Override
	public String getAddress() {
		return address;
	}

	@Override
	public Handler<Message<I,O>> getHandler() {
		return handler;
	}

	@Override
	public String getId(){
		return this.id;
	}
	
	private String genarateId(){
		return UUID.randomUUID().toString();
	}
}
