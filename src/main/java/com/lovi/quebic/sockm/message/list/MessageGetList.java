package com.lovi.quebic.sockm.message.list;

import java.util.ArrayList;
import java.util.List;

public class MessageGetList implements AutoCloseable {

	private String senderLauncherId_GET_LIST;
	private String contextKey_GET_LIST;
	private List<?> contextDataList_GET_LIST = new ArrayList<>();

	public String getSenderLauncherId_GET_LIST() {
		return senderLauncherId_GET_LIST;
	}

	public void setSenderLauncherId_GET_LIST(String senderLauncherId_GET_LIST) {
		this.senderLauncherId_GET_LIST = senderLauncherId_GET_LIST;
	}

	public String getContextKey_GET_LIST() {
		return contextKey_GET_LIST;
	}

	public void setContextKey_GET_LIST(String contextKey_GET_LIST) {
		this.contextKey_GET_LIST = contextKey_GET_LIST;
	}

	public List<?> getContextDataList_GET_LIST() {
		return contextDataList_GET_LIST;
	}

	public void setContextDataList_GET_LIST(
			List<?> contextDataList_GET_LIST) {
		this.contextDataList_GET_LIST = contextDataList_GET_LIST;
	}

	@Override
	public String toString() {
		return "MessageGetList [senderLauncherId_GET_LIST="
				+ senderLauncherId_GET_LIST + ", contextKey_GET_LIST="
				+ contextKey_GET_LIST + ", contextDataList_GET_LIST="
				+ contextDataList_GET_LIST + "]";
	}

	@Override
	public void close() throws Exception {

	}
}
