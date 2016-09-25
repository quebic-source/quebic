package com.lovi.quebic.sockm.message.map;

public class MessageRemoveMap implements AutoCloseable{

	private String senderLauncherId_REMOVE_MAP;
	private String contextKey_REMOVE_MAP;
	private Object mapKey_REMOVE_MAP;
	
	public String getSenderLauncherId_REMOVE_MAP() {
		return senderLauncherId_REMOVE_MAP;
	}
	public void setSenderLauncherId_REMOVE_MAP(String senderLauncherId_REMOVE_MAP) {
		this.senderLauncherId_REMOVE_MAP = senderLauncherId_REMOVE_MAP;
	}
	public String getContextKey_REMOVE_MAP() {
		return contextKey_REMOVE_MAP;
	}
	public void setContextKey_REMOVE_MAP(String contextKey_REMOVE_MAP) {
		this.contextKey_REMOVE_MAP = contextKey_REMOVE_MAP;
	}
	public Object getMapKey_REMOVE_MAP() {
		return mapKey_REMOVE_MAP;
	}
	public void setMapKey_REMOVE_MAP(Object mapKey_REMOVE_MAP) {
		this.mapKey_REMOVE_MAP = mapKey_REMOVE_MAP;
	}
	@Override
	public String toString() {
		return "MessageRemoveMap [senderLauncherId_REMOVE_MAP="
				+ senderLauncherId_REMOVE_MAP + ", contextKey_REMOVE_MAP="
				+ contextKey_REMOVE_MAP + ", mapKey_REMOVE_MAP="
				+ mapKey_REMOVE_MAP + "]";
	}
	@Override
	public void close() throws Exception {
	}
	
}
