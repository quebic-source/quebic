package com.lovi.quebic.web.impl;

import java.util.HashMap;
import java.util.Map;

import com.lovi.quebic.web.ApplicationContextData;

public class ApplicationContextDataImpl implements ApplicationContextData {

	private Map<String, Object> parms = new HashMap<>();
	
	@Override
	public void put(String key, Object value) {
		parms.put(key, value);
	}

	@SuppressWarnings("unchecked")
	@Override
	public <T> T get(String key) {
		return (T)parms.get(key);
	}

	@Override
	public Object remove(String key) {
		return parms.remove(key);
	}
	
	@Override
	public void clearAll() {
		parms.clear();
	}

	@Override
	public Map<String, Object> getParms() {
		return parms;
	}

	@Override
	public void setParms(Map<String, Object> parms) {
		this.parms = parms;
	}
}
