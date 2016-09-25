package com.lovi.quebic.sockm.message.map;

import java.util.HashMap;
import java.util.Map;

public class MessageGetMap implements AutoCloseable{
	
	private String senderLauncherId_GET_MAP;
	private String contextKey_GET_MAP;
	private Map<?, ?> contextDataMap_GET_MAP = new HashMap<>();
	
	public String getSenderLauncherId_GET_MAP() {
		return senderLauncherId_GET_MAP;
	}
	public void setSenderLauncherId_GET_MAP(String senderLauncherId_GET_MAP) {
		this.senderLauncherId_GET_MAP = senderLauncherId_GET_MAP;
	}
	public String getContextKey_GET_MAP() {
		return contextKey_GET_MAP;
	}
	public void setContextKey_GET_MAP(String contextKey_GET_MAP) {
		this.contextKey_GET_MAP = contextKey_GET_MAP;
	}
	public Map<?, ?> getContextDataMap_GET_MAP() {
		return contextDataMap_GET_MAP;
	}
	public void setContextDataMap_GET_MAP(Map<?, ?> contextDataMap_GET_MAP) {
		this.contextDataMap_GET_MAP = contextDataMap_GET_MAP;
	}
	@Override
	public String toString() {
		return "MessageGetMap [senderLauncherId_GET_MAP="
				+ senderLauncherId_GET_MAP + ", contextKey_GET_MAP="
				+ contextKey_GET_MAP + ", contextDataMap_GET_MAP="
				+ contextDataMap_GET_MAP + "]";
	}
	@Override
	public void close() throws Exception {
		
	}
}
