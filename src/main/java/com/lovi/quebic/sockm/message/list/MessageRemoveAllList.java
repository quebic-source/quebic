package com.lovi.quebic.sockm.message.list;

import java.util.Collection;

public class MessageRemoveAllList implements AutoCloseable {

	private String senderLauncherId_REMOVEALL_LIST;
	private String contextKey_REMOVEALL_LIST;
	private Collection<?> collection_REMOVEALL_LIST;

	public String getSenderLauncherId_REMOVEALL_LIST() {
		return senderLauncherId_REMOVEALL_LIST;
	}

	public void setSenderLauncherId_REMOVEALL_LIST(
			String senderLauncherId_REMOVEALL_LIST) {
		this.senderLauncherId_REMOVEALL_LIST = senderLauncherId_REMOVEALL_LIST;
	}

	public String getContextKey_REMOVEALL_LIST() {
		return contextKey_REMOVEALL_LIST;
	}

	public void setContextKey_REMOVEALL_LIST(String contextKey_REMOVEALL_LIST) {
		this.contextKey_REMOVEALL_LIST = contextKey_REMOVEALL_LIST;
	}

	public Collection<?> getCollection_REMOVEALL_LIST() {
		return collection_REMOVEALL_LIST;
	}

	public void setCollection_REMOVEALL_LIST(
			Collection<?> collection_REMOVEALL_LIST) {
		this.collection_REMOVEALL_LIST = collection_REMOVEALL_LIST;
	}
	
	@Override
	public String toString() {
		return "MessageRemoveAllList [senderLauncherId_REMOVEALL_LIST="
				+ senderLauncherId_REMOVEALL_LIST
				+ ", contextKey_REMOVEALL_LIST=" + contextKey_REMOVEALL_LIST
				+ ", collection_REMOVEALL_LIST=" + collection_REMOVEALL_LIST
				+ "]";
	}

	@Override
	public void close() throws Exception {
	}

}
