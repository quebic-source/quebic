package com.lovi.quebic.web;

/**
 * ViewAttribute is used to share data between controllers and views
 * @author Tharanga Thennakoon
 *
 */
public interface ViewAttribute {

	/**
	 * Put context data with key
	 * @param key
	 * @param object
	 */
	void put(String key,Object object);
	
	/**
	 * Get context data for the key
	 * @param key
	 * @return
	 */
	Object get(String key);
	
	/**
	 * Load template name
	 * @param name
	 */
	void loadView(String name);
}
