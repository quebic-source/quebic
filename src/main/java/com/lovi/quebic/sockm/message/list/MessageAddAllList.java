package com.lovi.quebic.sockm.message.list;

import java.util.Collection;

public class MessageAddAllList implements AutoCloseable {

	private String senderLauncherId_ADDALL_LIST;
	private String contextKey_ADDALL_LIST;
	private Collection<?> collection_ADDALL_LIST;
	private Integer listIndex_ADDALL_LIST = -1;// default

	public String getSenderLauncherId_ADDALL_LIST() {
		return senderLauncherId_ADDALL_LIST;
	}

	public void setSenderLauncherId_ADDALL_LIST(
			String senderLauncherId_ADDALL_LIST) {
		this.senderLauncherId_ADDALL_LIST = senderLauncherId_ADDALL_LIST;
	}

	public String getContextKey_ADDALL_LIST() {
		return contextKey_ADDALL_LIST;
	}

	public void setContextKey_ADDALL_LIST(String contextKey_ADDALL_LIST) {
		this.contextKey_ADDALL_LIST = contextKey_ADDALL_LIST;
	}

	public Collection<?> getCollection_ADDALL_LIST() {
		return collection_ADDALL_LIST;
	}

	public void setCollection_ADDALL_LIST(Collection<?> collection_ADDALL_LIST) {
		this.collection_ADDALL_LIST = collection_ADDALL_LIST;
	}

	public Integer getListIndex_ADDALL_LIST() {
		return listIndex_ADDALL_LIST;
	}

	public void setListIndex_ADDALL_LIST(Integer listIndex_ADDALL_LIST) {
		this.listIndex_ADDALL_LIST = listIndex_ADDALL_LIST;
	}
	
	@Override
	public String toString() {
		return "MessageAddAllList [senderLauncherId_ADDALL_LIST="
				+ senderLauncherId_ADDALL_LIST + ", contextKey_ADDALL_LIST="
				+ contextKey_ADDALL_LIST + ", collection_ADDALL_LIST="
				+ collection_ADDALL_LIST + ", listIndex_ADDALL_LIST="
				+ listIndex_ADDALL_LIST + "]";
	}

	@Override
	public void close() throws Exception {
	}

}
