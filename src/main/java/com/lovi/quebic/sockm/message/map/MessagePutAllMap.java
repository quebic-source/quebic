package com.lovi.quebic.sockm.message.map;

import java.util.HashMap;
import java.util.Map;

public class MessagePutAllMap implements AutoCloseable{
	
	private String senderLauncherId_PUTALL_MAP;
	private String contextKey_PUTALL_MAP;
	private Map<?, ?> map_PUTALL_MAP = new HashMap<>();
	
	public String getSenderLauncherId_PUTALL_MAP() {
		return senderLauncherId_PUTALL_MAP;
	}

	public void setSenderLauncherId_PUTALL_MAP(String senderLauncherId_PUTALL_MAP) {
		this.senderLauncherId_PUTALL_MAP = senderLauncherId_PUTALL_MAP;
	}

	public String getContextKey_PUTALL_MAP() {
		return contextKey_PUTALL_MAP;
	}

	public void setContextKey_PUTALL_MAP(String contextKey_PUTALL_MAP) {
		this.contextKey_PUTALL_MAP = contextKey_PUTALL_MAP;
	}

	public Map<?, ?> getMap_PUTALL_MAP() {
		return map_PUTALL_MAP;
	}

	public void setMap_PUTALL_MAP(Map<?, ?> map_PUTALL_MAP) {
		this.map_PUTALL_MAP = map_PUTALL_MAP;
	}
	
	@Override
	public String toString() {
		return "MessagePutAllMap [senderLauncherId_PUTALL_MAP="
				+ senderLauncherId_PUTALL_MAP + ", contextKey_PUTALL_MAP="
				+ contextKey_PUTALL_MAP + ", map_PUTALL_MAP=" + map_PUTALL_MAP
				+ "]";
	}

	@Override
	public void close() throws Exception {
		
	}
}
