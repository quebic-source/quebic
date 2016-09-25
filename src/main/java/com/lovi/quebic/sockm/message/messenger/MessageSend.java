package com.lovi.quebic.sockm.message.messenger;

public class MessageSend implements AutoCloseable {

	private String senderLauncherId_SEND;
	private String subscriberId_SEND;
	private Object value_SEND;
	
	public String getSenderLauncherId_SEND() {
		return senderLauncherId_SEND;
	}

	public void setSenderLauncherId_SEND(String senderLauncherId_SEND) {
		this.senderLauncherId_SEND = senderLauncherId_SEND;
	}

	public String getSubscriberId_SEND() {
		return subscriberId_SEND;
	}
	
	public void setSubscriberId_SEND(String subscriberId_SEND) {
		this.subscriberId_SEND = subscriberId_SEND;
	}

	public Object getValue_SEND() {
		return value_SEND;
	}

	public void setValue_SEND(Object value_SEND) {
		this.value_SEND = value_SEND;
	}

	@Override
	public void close() throws Exception {
	}

}
