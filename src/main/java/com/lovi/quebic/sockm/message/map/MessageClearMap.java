package com.lovi.quebic.sockm.message.map;

public class MessageClearMap implements AutoCloseable {

	private String senderLauncherId_CLEAR_MAP;
	private String contextKey_CLEAR_MAP;

	public String getSenderLauncherId_CLEAR_MAP() {
		return senderLauncherId_CLEAR_MAP;
	}

	public void setSenderLauncherId_CLEAR_MAP(String senderLauncherId_CLEAR_MAP) {
		this.senderLauncherId_CLEAR_MAP = senderLauncherId_CLEAR_MAP;
	}

	public String getContextKey_CLEAR_MAP() {
		return contextKey_CLEAR_MAP;
	}

	public void setContextKey_CLEAR_MAP(String contextKey_CLEAR_MAP) {
		this.contextKey_CLEAR_MAP = contextKey_CLEAR_MAP;
	}

	@Override
	public String toString() {
		return "MessageClearMap [senderLauncherId_CLEAR_MAP="
				+ senderLauncherId_CLEAR_MAP + ", contextKey_CLEAR_MAP="
				+ contextKey_CLEAR_MAP + "]";
	}

	@Override
	public void close() throws Exception {
	}

}
