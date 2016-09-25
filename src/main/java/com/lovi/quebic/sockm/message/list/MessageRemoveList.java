package com.lovi.quebic.sockm.message.list;

public class MessageRemoveList implements AutoCloseable {

	private String senderLauncherId_REMOVE_LIST;
	private String contextKey_REMOVE_LIST;
	private Object listValue_REMOVE_LIST;

	public String getSenderLauncherId_REMOVE_LIST() {
		return senderLauncherId_REMOVE_LIST;
	}

	public void setSenderLauncherId_REMOVE_LIST(
			String senderLauncherId_REMOVE_LIST) {
		this.senderLauncherId_REMOVE_LIST = senderLauncherId_REMOVE_LIST;
	}

	public String getContextKey_REMOVE_LIST() {
		return contextKey_REMOVE_LIST;
	}

	public void setContextKey_REMOVE_LIST(String contextKey_REMOVE_LIST) {
		this.contextKey_REMOVE_LIST = contextKey_REMOVE_LIST;
	}

	public Object getListValue_REMOVE_LIST() {
		return listValue_REMOVE_LIST;
	}

	public void setListValue_REMOVE_LIST(Object listValue_REMOVE_LIST) {
		this.listValue_REMOVE_LIST = listValue_REMOVE_LIST;
	}

	@Override
	public String toString() {
		return "MessageRemoveList [senderLauncherId_REMOVE_LIST="
				+ senderLauncherId_REMOVE_LIST + ", contextKey_REMOVE_LIST="
				+ contextKey_REMOVE_LIST + ", listValue_REMOVE_LIST="
				+ listValue_REMOVE_LIST + "]";
	}

	@Override
	public void close() throws Exception {
	}

}
