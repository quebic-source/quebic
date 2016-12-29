package com.lovi.quebic.web.impl;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;
import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.cluster.message.CoreMessage;
import com.lovi.quebic.cluster.message.MessageKey;
import com.lovi.quebic.cluster.option.MulticastGroup;
import com.lovi.quebic.exception.ErrorMessage;
import com.lovi.quebic.web.SessionStore;

public class ClusterSessionStoreImpl implements SessionStore {
	
	final static Logger logger = LoggerFactory.getLogger(ClusterSessionStoreImpl.class);
	
	private Map<String, Map<String, Object>> users = new HashMap<>();
	
	private final String multicastAddress;
	private final int multicastPort;
	
	public ClusterSessionStoreImpl(MulticastGroup multicastGroup) {
		multicastAddress = multicastGroup.getMulticastAddress();
		multicastPort = multicastGroup.getMulticastPort();
	}
	
	@Override
	public String addNewUser(){
		
		String sessionId = genarateSessionId();
		while(users.containsKey(sessionId)){
			sessionId = genarateSessionId();
		}
		
		users.put(sessionId, new HashMap<>());
		
		publishAddNewUser(sessionId);
		
		return sessionId;
	}
	
	@Override
	public void removeUser(String sessionId){
		
		if(!users.containsKey(sessionId)){
			logger.error(ErrorMessage.SESSION_USER_ALREADY_REMOVED.getMessage());
			return;
		}
			
		users.remove(sessionId);
		publishRemoveUser(sessionId);
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
	
	private void publishAddNewUser(String sessionId){
		try {

			CoreMessage coreMessage = new CoreMessage();
			coreMessage.setKey(MessageKey.SESSION_NEW_USER);
			
			Map<String, Map<String,Object>> sessionParm = new HashMap<>();
			sessionParm.put(sessionId, new HashMap<>());
			
			coreMessage.setSessionParm(sessionParm);
			
			sendMulticastMessage(coreMessage);
			
		} catch (Exception e) {
			logger.error("publishAddNewUser " + e.getMessage());
		}
	}
	
	private void publishRemoveUser(String sessionId){
		try {

			CoreMessage coreMessage = new CoreMessage();
			coreMessage.setKey(MessageKey.SESSION_REMOVE_USER);
			
			Map<String, Map<String,Object>> sessionParm = new HashMap<>();
			sessionParm.put(sessionId, new HashMap<>());
			
			coreMessage.setSessionParm(sessionParm);
			
			sendMulticastMessage(coreMessage);
			
		} catch (Exception e) {
			logger.error("publishRemoveUser " + e.getMessage());
		}
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
	public String toString() {
		return "SessionStore [users=" + users + "]";
	}
	
}
