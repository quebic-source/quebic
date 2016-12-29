package com.lovi.quebic.web.impl;

import java.util.Map;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.exception.ErrorMessage;
import com.lovi.quebic.web.Session;
import com.lovi.quebic.web.SessionStore;

public class SessionImpl implements Session{
	
	final static Logger logger = LoggerFactory.getLogger(SessionImpl.class);
	
	private SessionStore sessionStore;
	private String sessionId;
	
	public SessionImpl(SessionStore sessionStore, String sessionId) {
		this.sessionStore = sessionStore;
		this.sessionId = sessionId;
	}
	
	@Override
	public void put(String key, Object value){
		sessionStore.getUsers().get(sessionId).put(key, value);
	}
	
	@Override
	public Object get(String key){
		
		if(!sessionStore.getUsers().containsKey(sessionId)){
			logger.error(ErrorMessage.SESSION_USER_NOT_FOUND.getMessage());
			return null;
		}
		
		return sessionStore.getUsers().get(sessionId).get(key);
	}
	
	@Override
	public Object remove(String key){
		
		if(!sessionStore.getUsers().containsKey(sessionId)){
			logger.error(ErrorMessage.SESSION_USER_NOT_FOUND.getMessage());
			return null;
		}
		
		if(!sessionStore.getUsers().get(sessionId).containsKey(key)){
			logger.error(ErrorMessage.SESSION_PARM_NOT_FOUND.getMessage() + key);
			return null;
		}
		
		return sessionStore.getUsers().get(sessionId).remove(key);
	}
	
	@Override
	public void clearAll() {
		
		if(!sessionStore.getUsers().containsKey(sessionId)){
			logger.error(ErrorMessage.SESSION_USER_NOT_FOUND.getMessage());
			return;
		}
		
		sessionStore.getUsers().get(sessionId).clear();
		
	}

	@Override
	public void releaseMe() {
		sessionStore.removeUser(sessionId);
	}

	@Override
	public Map<String, Object> getParms() {
		return sessionStore.getUsers().get(sessionId);
	}

	
}
