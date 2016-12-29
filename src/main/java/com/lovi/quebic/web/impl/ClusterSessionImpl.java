package com.lovi.quebic.web.impl;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.HashMap;
import java.util.Map;
import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.cluster.message.CoreMessage;
import com.lovi.quebic.cluster.message.MessageKey;
import com.lovi.quebic.cluster.option.MulticastGroup;
import com.lovi.quebic.exception.ErrorMessage;
import com.lovi.quebic.web.Session;
import com.lovi.quebic.web.SessionStore;

public class ClusterSessionImpl implements Session{

	final static Logger logger = LoggerFactory.getLogger(ClusterSessionImpl.class);
	
	private SessionStore sessionStore;
	private String sessionId;
	
	private final String multicastAddress;
	private final int multicastPort;
	
	public ClusterSessionImpl(SessionStore sessionStore, String sessionId, MulticastGroup multicastGroup) {
		this.sessionStore = sessionStore;
		this.sessionId = sessionId;
		
		multicastAddress = multicastGroup.getMulticastAddress();
		multicastPort = multicastGroup.getMulticastPort();
	}
	
	@Override
	public void put(String key, Object value){
		sessionStore.getUsers().get(sessionId).put(key, value);
		publishSessionPut(key, value);
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
		
		Object removedObject = sessionStore.getUsers().get(sessionId).remove(key);
		publishSessionRemove(key);
		return removedObject;
	}
	
	@Override
	public void clearAll() {
		
		if(!sessionStore.getUsers().containsKey(sessionId)){
			logger.error(ErrorMessage.SESSION_USER_NOT_FOUND.getMessage());
			return;
		}
		
		sessionStore.getUsers().get(sessionId).clear();
		publishSessionClearAll();
		
	}
	
	private void publishSessionPut(String key, Object value){
		try {

			CoreMessage coreMessage = new CoreMessage();
			coreMessage.setKey(MessageKey.SESSION_PARM_PUT);
			
			Map<String, Map<String,Object>> sessionParm = new HashMap<>();
			Map<String, Object> sessionPair = new HashMap<>();
			sessionPair.put(key, value);
			sessionParm.put(sessionId, sessionPair);
			
			coreMessage.setSessionParm(sessionParm);
			
			sendMulticastMessage(coreMessage);
			
		} catch (Exception e) {
			logger.error("publishSessionChange " + e.getMessage());
		}
	}
	
	private void publishSessionRemove(String key){
		try {

			CoreMessage coreMessage = new CoreMessage();
			coreMessage.setKey(MessageKey.SESSION_PARM_REMOVE);
			
			Map<String, Map<String,Object>> sessionParm = new HashMap<>();
			Map<String, Object> sessionPair = new HashMap<>();
			sessionPair.put(key, null);
			sessionParm.put(sessionId, sessionPair);
			
			coreMessage.setSessionParm(sessionParm);
			
			sendMulticastMessage(coreMessage);
			
		} catch (Exception e) {
			logger.error("publishSessionChange " + e.getMessage());
		}
	}

	private void publishSessionClearAll(){
		try {

			CoreMessage coreMessage = new CoreMessage();
			coreMessage.setKey(MessageKey.SESSION_PARMS_CLEAR_ALL);
			
			Map<String, Map<String,Object>> sessionParm = new HashMap<>();
			sessionParm.put(sessionId, new HashMap<>());
			
			coreMessage.setSessionParm(sessionParm);
			
			sendMulticastMessage(coreMessage);
			
		} catch (Exception e) {
			logger.error("publishSessionClearAll " + e.getMessage());
		}
	}
	
	@Override
	public void releaseMe() {
		sessionStore.removeUser(sessionId);
	}
	
	private void sendMulticastMessage(CoreMessage coreMessage)throws Exception{
		
		DatagramSocket udpSocket = new DatagramSocket();

		InetAddress mcIPAddress = InetAddress.getByName(multicastAddress);

		ObjectMapper mapper = new ObjectMapper();
		byte[] messageBytes = mapper.writeValueAsString(
				coreMessage).getBytes();

		DatagramPacket packet = new DatagramPacket(
				messageBytes, messageBytes.length);
		packet.setAddress(mcIPAddress);
		packet.setPort(multicastPort);
		udpSocket.send(packet);
		udpSocket.close();
		
	}

	@Override
	public Map<String, Object> getParms() {
		return sessionStore.getUsers().get(sessionId);
	}
}
