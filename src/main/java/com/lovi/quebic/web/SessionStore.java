package com.lovi.quebic.web;

import java.util.Map;

public interface SessionStore {

	final static String sessionIdKey = "puppy-io.sessionid";
	
	String addNewUser();
	
	void removeUser(String sessionId);

	boolean checkSessionIdExists(String sessionId);
	
	Map<String, Map<String, Object>> getUsers();
	
	void setUsers(Map<String, Map<String, Object>> users);

}
