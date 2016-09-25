package com.lovi.quebic.sockm.message.map;

public class MessagePutMap implements AutoCloseable{

	private String senderLauncherId_PUT_MAP;
	private String contextKey_PUT_MAP;
	private Object mapKey_PUT_MAP;
	private Object mapValue_PUT_MAP;
	
	public String getSenderLauncherId_PUT_MAP() {
		return senderLauncherId_PUT_MAP;
	}
	public void setSenderLauncherId_PUT_MAP(String senderLauncherId_PUT_MAP) {
		this.senderLauncherId_PUT_MAP = senderLauncherId_PUT_MAP;
	}
	public String getContextKey_PUT_MAP() {
		return contextKey_PUT_MAP;
	}
	public void setContextKey_PUT_MAP(String contextKey_PUT_MAP) {
		this.contextKey_PUT_MAP = contextKey_PUT_MAP;
	}
	public Object getMapKey_PUT_MAP() {
		return mapKey_PUT_MAP;
	}
	public void setMapKey_PUT_MAP(Object mapKey_PUT_MAP) {
		this.mapKey_PUT_MAP = mapKey_PUT_MAP;
	}
	public Object getMapValue_PUT_MAP() {
		return mapValue_PUT_MAP;
	}
	public void setMapValue_PUT_MAP(Object mapValue_PUT_MAP) {
		this.mapValue_PUT_MAP = mapValue_PUT_MAP;
	}
	@Override
	public String toString() {
		return "MessagePutMap [senderLauncherId_PUT_MAP="
				+ senderLauncherId_PUT_MAP + ", contextKey_PUT_MAP="
				+ contextKey_PUT_MAP + ", mapKey_PUT_MAP=" + mapKey_PUT_MAP
				+ ", mapValue_PUT_MAP=" + mapValue_PUT_MAP + "]";
	}
	@Override
	public void close() throws Exception {
	}
	
}
