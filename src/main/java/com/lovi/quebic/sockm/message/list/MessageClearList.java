package com.lovi.quebic.sockm.message.list;

public class MessageClearList implements AutoCloseable {

	private String senderLauncherId_CLEAR_LIST;
	private String contextKey_CLEAR_LIST;

	public String getSenderLauncherId_CLEAR_LIST() {
		return senderLauncherId_CLEAR_LIST;
	}

	public void setSenderLauncherId_CLEAR_LIST(
			String senderLauncherId_CLEAR_LIST) {
		this.senderLauncherId_CLEAR_LIST = senderLauncherId_CLEAR_LIST;
	}

	public String getContextKey_CLEAR_LIST() {
		return contextKey_CLEAR_LIST;
	}

	public void setContextKey_CLEAR_LIST(String contextKey_CLEAR_LIST) {
		this.contextKey_CLEAR_LIST = contextKey_CLEAR_LIST;
	}
	
	@Override
	public String toString() {
		return "MessageClearList [senderLauncherId_CLEAR_LIST="
				+ senderLauncherId_CLEAR_LIST + ", contextKey_CLEAR_LIST="
				+ contextKey_CLEAR_LIST + "]";
	}

	@Override
	public void close() throws Exception {
	}

}
