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
import com.lovi.quebic.web.ApplicationContextData;

public class ClusterApplicationContextDataImpl implements ApplicationContextData {

	final static Logger logger = LoggerFactory.getLogger(ClusterApplicationContextDataImpl.class);
	private Map<String, Object> parms = new HashMap<>();
	
	private final String multicastAddress;
	private final int multicastPort;
	
	public ClusterApplicationContextDataImpl(MulticastGroup multicastGroup) {
		multicastAddress = multicastGroup.getMulticastAddress();
		multicastPort = multicastGroup.getMulticastPort();
	}
	
	@Override
	public void put(String key, Object value) {
		parms.put(key, value);
		publishApplicationParmPut(key, value);
	}

	@SuppressWarnings("unchecked")
	@Override
	public <T> T get(String key) {
		return (T) parms.get(key);
	}

	@Override
	public Object remove(String key) {
		Object removedObject = parms.remove(key);
		publishApplicationParmRemove(key);
		return removedObject;
	}
	
	@Override
	public void clearAll() {
		parms.clear();
		publishApplicationParmClearAll();
	}

	@Override
	public Map<String, Object> getParms() {
		return parms;
	}
	
	private void publishApplicationParmPut(String key, Object value){
		try {

			CoreMessage coreMessage = new CoreMessage();
			coreMessage.setKey(MessageKey.APPLICATION_PARM_PUT);
			
			Map<String, Object> parmPair = new HashMap<>();
			parmPair.put(key, value);
			
			coreMessage.setApplicationContextParms(parmPair);
			
			sendMulticastMessage(coreMessage);
			
		} catch (Exception e) {
			logger.error("publishApplicationParmPut " + e.getMessage());
		}
	}
	
	private void publishApplicationParmRemove(String key){
		try {

			CoreMessage coreMessage = new CoreMessage();
			coreMessage.setKey(MessageKey.APPLICATION_PARM_REMOVE);
			
			Map<String, Object> parmPair = new HashMap<>();
			parmPair.put(key, null);
			
			coreMessage.setApplicationContextParms(parmPair);
			
			sendMulticastMessage(coreMessage);
			
		} catch (Exception e) {
			logger.error("publishApplicationParmRemove " + e.getMessage());
		}
	}

	private void publishApplicationParmClearAll(){
		try {

			CoreMessage coreMessage = new CoreMessage();
			coreMessage.setKey(MessageKey.APPLICATION_PARMS_CLEAR_ALL);
			
			sendMulticastMessage(coreMessage);
			
		} catch (Exception e) {
			logger.error("publishApplicationParmClearAll " + e.getMessage());
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
	public void setParms(Map<String, Object> parms) {
		this.parms = parms;
	}

}
