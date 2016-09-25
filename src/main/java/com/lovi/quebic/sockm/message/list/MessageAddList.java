package com.lovi.quebic.sockm.message.list;

public class MessageAddList implements AutoCloseable {

	private String senderLauncherId_ADD_LIST;
	private String contextKey_ADD_LIST;
	private Object listValue_ADD_LIST;
	private Integer listIndex_ADD_LIST = -1;//default
	
	public String getSenderLauncherId_ADD_LIST() {
		return senderLauncherId_ADD_LIST;
	}

	public void setSenderLauncherId_ADD_LIST(String senderLauncherId_ADD_LIST) {
		this.senderLauncherId_ADD_LIST = senderLauncherId_ADD_LIST;
	}

	public String getContextKey_ADD_LIST() {
		return contextKey_ADD_LIST;
	}

	public void setContextKey_ADD_LIST(String contextKey_ADD_LIST) {
		this.contextKey_ADD_LIST = contextKey_ADD_LIST;
	}

	public Object getListValue_ADD_LIST() {
		return listValue_ADD_LIST;
	}

	public void setListValue_ADD_LIST(Object listValue_ADD_LIST) {
		this.listValue_ADD_LIST = listValue_ADD_LIST;
	}
	
	public Integer getListIndex_ADD_LIST() {
		return listIndex_ADD_LIST;
	}

	public void setListIndex_ADD_LIST(Integer listIndex_ADD_LIST) {
		this.listIndex_ADD_LIST = listIndex_ADD_LIST;
	}

	@Override
	public String toString() {
		return "MessageAddList [senderLauncherId_ADD_LIST="
				+ senderLauncherId_ADD_LIST + ", contextKey_ADD_LIST="
				+ contextKey_ADD_LIST + ", listValue_ADD_LIST="
				+ listValue_ADD_LIST + ", listIndex_ADD_LIST="
				+ listIndex_ADD_LIST + "]";
	}

	@Override
	public void close() throws Exception {
	}

}
