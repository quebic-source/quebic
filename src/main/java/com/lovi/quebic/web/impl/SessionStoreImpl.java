package com.lovi.quebic.web.impl;

import java.util.HashMap;
import java.util.Map;
import java.util.UUID;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.exception.ErrorMessage;
import com.lovi.quebic.web.SessionStore;

public class SessionStoreImpl implements SessionStore {
	
	final static Logger logger = LoggerFactory.getLogger(SessionStoreImpl.class);
	
	private Map<String, Map<String, Object>> users = new HashMap<>();
	
	@Override
	public String addNewUser(){
		
		String sessionId = genarateSessionId();
		while(users.containsKey(sessionId)){
			sessionId = genarateSessionId();
		}
		
		users.put(sessionId, new HashMap<>());
		return sessionId;
	}
	
	@Override
	public void removeUser(String sessionId){
		
		if(!users.containsKey(sessionId)){
			logger.error(ErrorMessage.SESSION_USER_ALREADY_REMOVED.getMessage());
			return;
		}
		
		users.remove(sessionId);
	}
	
	@Override
	public boolean checkSessionIdExists(String sessionId){
		if(users.containsKey(sessionId))
			return true;
		else
			return false;
	}

	@Override
	public Map<String, Map<String, Object>> getUsers() {
		return users;
	}
	
	@Override
	public void setUsers(Map<String, Map<String, Object>> users) {
		this.users = users;
	}
	
	private String genarateSessionId(){
		String sessionId = UUID.randomUUID().toString();
		return sessionId;
	}

	@Override
	public String toString() {
		return "SessionStore [users=" + users + "]";
	}
	
}
