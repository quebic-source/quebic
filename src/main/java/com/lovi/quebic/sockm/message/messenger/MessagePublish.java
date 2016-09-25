package com.lovi.quebic.sockm.message.messenger;

public class MessagePublish implements AutoCloseable {

	private String senderLauncherId_PUBLISH;
	private String address_PUBLISH;
	private Object value_PUBLISH;

	public String getSenderLauncherId_PUBLISH() {
		return senderLauncherId_PUBLISH;
	}

	public void setSenderLauncherId_PUBLISH(String senderLauncherId_PUBLISH) {
		this.senderLauncherId_PUBLISH = senderLauncherId_PUBLISH;
	}

	public String getAddress_PUBLISH() {
		return address_PUBLISH;
	}

	public void setAddress_PUBLISH(String address_PUBLISH) {
		this.address_PUBLISH = address_PUBLISH;
	}

	public Object getValue_PUBLISH() {
		return value_PUBLISH;
	}

	public void setValue_PUBLISH(Object value_PUBLISH) {
		this.value_PUBLISH = value_PUBLISH;
	}

	@Override
	public void close() throws Exception {
	}

}
