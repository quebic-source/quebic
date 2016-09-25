package com.lovi.quebic.sockm.message.messenger;

public class MessageSendReply {

	private Object responseValue;
	private String responseThrowableMessage = null;
	
	public MessageSendReply() {
	}

	public MessageSendReply(Object responseValue, String responseThrowableMessage) {
		super();
		this.responseValue = responseValue;
		this.responseThrowableMessage = responseThrowableMessage;
	}

	public Object getResponseValue() {
		return responseValue;
	}

	public void setResponseValue(Object responseValue) {
		this.responseValue = responseValue;
	}

	public String getResponseThrowableMessage() {
		return responseThrowableMessage;
	}

	public void setResponseThrowableMessage(String responseThrowableMessage) {
		this.responseThrowableMessage = responseThrowableMessage;
	}
	
	
	
	
}
