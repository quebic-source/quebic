package com.lovi.quebic.web;

import java.util.Map;

public interface Session {

	void put(String key, Object value);
	
	Object get(String key);

	Object remove(String key);
	
	void clearAll();

	void releaseMe();
	
	Map<String, Object> getParms();
}
