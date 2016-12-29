package com.lovi.quebic.sockm.context.impl;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.HashMap;
import java.util.Map;

import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.sockm.config.MulticastNetworkGroup;
import com.lovi.quebic.sockm.context.ContextMapData;
import com.lovi.quebic.sockm.launcher.SockMLauncher;
import com.lovi.quebic.sockm.log.info.InfoMessage;
import com.lovi.quebic.sockm.message.map.MessageGetMap;
import com.lovi.quebic.sockm.util.SharedMap;

public class ContextMapDataImpl implements ContextMapData{

	private final Logger logger = LoggerFactory.getLogger(ContextMapData.class);
	
	private SockMLauncher launcher;
	private MulticastNetworkGroup multicastGroup;
	private Map<String, Map<?, ?>> dataMap = new HashMap<>();
	
	public ContextMapDataImpl(SockMLauncher launcher) {
		this.launcher = launcher;
		multicastGroup = launcher.getMulticastGroup();
	}
	
	@Override
	public <K, V> Map<K, V> getMap(String key){
		
		if(dataMap.containsKey(key)){
			@SuppressWarnings("unchecked")
			Map<K, V> map = (Map<K, V>) dataMap.get(key);
			return map;
		}
		else{
			Map<K, V> map = createtMapData(key, null, null);
			dataMap.put(key, map);
			return map;
		}
	}
	
	@Override
	public <K, V> Map<K, V> getMap(String key, Class<K> typeOfKey, Class<V> typeOfValue){
		
		if(dataMap.containsKey(key)){
			@SuppressWarnings("unchecked")
			Map<K, V> map = (Map<K, V>) dataMap.get(key);
			return map;
		}
		else{
			Map<K, V> map = createtMapData(key, typeOfKey, typeOfValue);
			dataMap.put(key, map);
			return map;
		}
	}
	
	@Override
	public void updateMapForPut(String contextKey, Object key, Object value){
		if(dataMap.containsKey(contextKey)){
			@SuppressWarnings("unchecked")
			SharedMap<Object, Object> map = (SharedMap<Object, Object>) dataMap.get(contextKey);
			map._put(key, value);
		}
	}
	
	@Override
	public void updateMapForPutAll(String contextKey, Map<?, ?> m){
		if(dataMap.containsKey(contextKey)){
			@SuppressWarnings("unchecked")
			SharedMap<Object, Object> map = (SharedMap<Object, Object>) dataMap.get(contextKey);
			map._putAll(m);
		}
	}

	@Override
	public void updateMapForRemove(String contextKey, Object key) {
		if(dataMap.containsKey(contextKey)){
			@SuppressWarnings("unchecked")
			SharedMap<Object, Object> map = (SharedMap<Object, Object>) dataMap.get(contextKey);
			map._remove(key);
		}
	}
	
	@Override
	public void updateMapForClear(String contextKey) {
		if(dataMap.containsKey(contextKey)){
			@SuppressWarnings("unchecked")
			SharedMap<Object, Object> map = (SharedMap<Object, Object>) dataMap.get(contextKey);
			map._clear();
		}
	}
	
	@Override
	public Map<String, Map<?, ?>> getDataMap(){
		return dataMap;
	}
	
	@SuppressWarnings("unchecked")
	private <K, V> Map<K, V> createtMapData(String contextKey, Class<K> typeOfKey, Class<V> typeOfValue){
		SharedMap<K, V> map = new SharedMap<>(launcher, contextKey , typeOfKey, typeOfValue);
		try {
			
			DatagramSocket udpSocket = new DatagramSocket();
			
			InetAddress mcIPAddress = InetAddress.getByName(multicastGroup.getMulticastAddress());
			
			MessageGetMap messageGetMap = new MessageGetMap();
			messageGetMap.setSenderLauncherId_GET_MAP(launcher.getLauncherId());
			messageGetMap.setContextKey_GET_MAP(contextKey);

			ObjectMapper mapper = new ObjectMapper();
			byte[] joinRequestMessageBytes = mapper.writeValueAsString(
					messageGetMap).getBytes();
			
			DatagramPacket requestPacket = new DatagramPacket(
					joinRequestMessageBytes, joinRequestMessageBytes.length);
			requestPacket.setAddress(mcIPAddress);
			requestPacket.setPort(multicastGroup.getMulticastPort());
			udpSocket.send(requestPacket);
			
			DatagramPacket replyPacket = new DatagramPacket(new byte[1024],
					1024);
			replyPacket.setAddress(mcIPAddress);
			replyPacket.setPort(multicastGroup.getMulticastPort());

			try {
				
				udpSocket.setSoTimeout(500);
				udpSocket.receive(replyPacket);

				String replyMsgStr = new String(replyPacket.getData(),
						replyPacket.getOffset(), replyPacket.getLength());

				messageGetMap = mapper.readValue(replyMsgStr, MessageGetMap.class);
				
				map._putAll((Map<? extends K, ? extends V>) messageGetMap.getContextDataMap_GET_MAP());

			} catch (Exception e) {
				logger.info(InfoMessage.CREATE_NEW_MAP);
				//logger.error(e.getMessage());
			}
			
			udpSocket.close();
			
		} catch (Exception e) {
			logger.error("createtMapData " + e.getMessage());
		}
		
		return map;
	}
}
