package com.lovi.quebic.sockm.message.messenger;

import com.lovi.quebic.async.Future;

public class Message<I,O>{

	private I requestMessage;
	private Future<O> responseFuture;
	
	public Message() {
	}
	
	public Message(I requestMessage) {
		this.requestMessage = requestMessage;
		responseFuture = Future.create();
	}
	
	public void setMessage(I requestMessage){
		this.requestMessage = requestMessage;
	}
	public I getMessage(){
		return requestMessage;
	}
	
	public void reply(O message){
		responseFuture.setResult(message);
	}
	
	public void replyFailure(Throwable failure){
		responseFuture.setFailure(failure);
	}
	
	public O getReply(){
		return responseFuture.getResult();
	}
	
	public Throwable getReplyFailure(){
		return responseFuture.getFailure();
	}
	
	public boolean isReplyFail(){
		return responseFuture.isFail();
	}

	@Override
	public String toString() {
		return "Message [requestMessage=" + requestMessage
				+ ", responseFuture=" + responseFuture + "]";
	}
	
}
