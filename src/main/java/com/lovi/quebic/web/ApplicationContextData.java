package com.lovi.quebic.web;

import java.util.Map;

public interface ApplicationContextData {

	void put(String key, Object value);
	
	<T> T get(String key);

	Object remove(String key);
	
	void clearAll();
	
	Map<String,Object> getParms();
	
	void setParms(Map<String,Object> parms);
}
