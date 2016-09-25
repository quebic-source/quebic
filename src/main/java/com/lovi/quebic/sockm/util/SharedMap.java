package com.lovi.quebic.sockm.util;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.sockm.config.MulticastNetworkGroup;
import com.lovi.quebic.sockm.launcher.SockMLauncher;
import com.lovi.quebic.sockm.message.map.MessageClearMap;
import com.lovi.quebic.sockm.message.map.MessagePutAllMap;
import com.lovi.quebic.sockm.message.map.MessagePutMap;
import com.lovi.quebic.sockm.message.map.MessageRemoveMap;

public final class SharedMap<K, V> implements Map<K, V>{

	private final Logger logger = LoggerFactory.getLogger(SharedMap.class);
	
	private Class<K> typeOfKey;
	private Class<V> typeOfValue;
	
	private Map<K, V> map = new HashMap<>();
	
	private SockMLauncher launcher;
	private String contextKey;
	private MulticastNetworkGroup multicastGroup;
	
	public SharedMap(SockMLauncher launcher, String contextKey, Class<K> typeOfKey, Class<V> typeOfValue) {
		this.launcher = launcher;
		this.multicastGroup = launcher.getMulticastGroup();
		this.contextKey = contextKey;
		this.typeOfKey = typeOfKey;
		this.typeOfValue = typeOfValue;
	}

	@Override
	public void clear() {
		
		MessageClearMap clearMap = new MessageClearMap();
		clearMap.setSenderLauncherId_CLEAR_MAP(launcher.getLauncherId());
		clearMap.setContextKey_CLEAR_MAP(contextKey);
		
		sendMessage(clearMap);
		
		map.clear();
	}

	@Override
	public boolean containsKey(Object key) {
		return map.containsKey(key);
	}

	@Override
	public boolean containsValue(Object value) {
		return map.containsValue(value);
	}

	@Override
	public Set<Entry<K, V>> entrySet() {
		return map.entrySet();
	}

	@Override
	public V get(Object key) {
		return map.get(key);
	}

	@Override
	public boolean isEmpty() {
		return map.isEmpty();
	}

	@Override
	public Set<K> keySet() {
		return map.keySet();
	}

	@Override
	public V put(K key, V value) {
		
		MessagePutMap putMap = new MessagePutMap();
		putMap.setSenderLauncherId_PUT_MAP(launcher.getLauncherId());
		putMap.setContextKey_PUT_MAP(contextKey);
		putMap.setMapKey_PUT_MAP(key);
		putMap.setMapValue_PUT_MAP(value);
		
		sendMessage(putMap);
		
		return map.put(key, value);
	}
	
	

	@SuppressWarnings("unchecked")
	@Override
	public void putAll(Map<? extends K, ? extends V> m) {
		
		MessagePutAllMap putAllMap = new MessagePutAllMap();
		putAllMap.setSenderLauncherId_PUTALL_MAP(launcher.getLauncherId());
		putAllMap.setContextKey_PUTALL_MAP(contextKey);
		putAllMap.setMap_PUTALL_MAP((Map<Object, Object>) m);
		
		sendMessage(putAllMap);
		
		map.putAll(m);
	}

	@Override
	public V remove(Object key) {
		
		MessageRemoveMap removeMap = new MessageRemoveMap();
		removeMap.setSenderLauncherId_REMOVE_MAP(launcher.getLauncherId());
		removeMap.setContextKey_REMOVE_MAP(contextKey);
		removeMap.setMapKey_REMOVE_MAP(key);
		
		sendMessage(removeMap);
		
		return map.remove(key);
	}
	
	public void _put(K key, V value){
		ObjectMapper objectMapper = new ObjectMapper();
		
		if(typeOfKey != null){
			try {
				
				if(!key.getClass().getName().equals(typeOfKey.getName())){
					String jsonStr = objectMapper.writeValueAsString(key);
					key = objectMapper.readValue(jsonStr, typeOfKey);
				}
				
			} catch (Exception e) {
				logger.error("_put process key " + e.getMessage());
			}

		}
		
		if(typeOfValue != null){
			try {
				
				if(!value.getClass().getName().equals(typeOfValue.getName())){
					String jsonStr = objectMapper.writeValueAsString(value);
					value = objectMapper.readValue(jsonStr, typeOfValue);
				}
				
			} catch (Exception e) {
				logger.error("_put process value " + e.getMessage());
			}

		}
		
		
		map.put(key, value);
	}
	
	public void _putAll(Map<? extends K, ? extends V> m) {
		
		ObjectMapper objectMapper = new ObjectMapper();
		
		for(Entry<? extends K, ? extends V> entry : m.entrySet()){
			
			K key = entry.getKey();
			V value = entry.getValue();
			
			if(typeOfKey != null){
				try {		
					
					if(!key.getClass().getName().equals(typeOfKey.getName())){
						String jsonStr = objectMapper.writeValueAsString(key);
						key = objectMapper.readValue(jsonStr, typeOfKey);
					}
					
				} catch (Exception e) {
					logger.error("_putAll process key " + e.getMessage());
				}
			}
			
			if(typeOfValue != null){
				try {
					
					if(!value.getClass().getName().equals(typeOfValue.getName())){
						String jsonStr = objectMapper.writeValueAsString(value);
						value = objectMapper.readValue(jsonStr, typeOfValue);
					}
					
				} catch (Exception e) {
					logger.error("_putAll process value " + e.getMessage());
				}
			}
			
			map.put(key, value);
		}
		
	}
	
	public void _remove(Object key) {
		map.remove(key);
	}
	
	public void _clear() {
		map.clear();
	}

	@Override
	public int size() {
		return map.size();
	}

	@Override
	public Collection<V> values() {
		return map.values();
	}
	
	@Override
	public String toString() {
		return "SharedMap [map=" + map + ", contextKey=" + contextKey + "]";
	}

	private <T> void sendMessage(T message){
		if(launcher.getDiscoveryOption().isEnableMulticastMethod()){
			try{
				DatagramSocket udpSocket = new DatagramSocket();
		
				InetAddress mcIPAddress = InetAddress.getByName(multicastGroup.getMulticastAddress());
		
				ObjectMapper mapper = new ObjectMapper();
				byte[] messageBytes = mapper.writeValueAsString(message).getBytes();
		
				DatagramPacket packet = new DatagramPacket(
						messageBytes, messageBytes.length);
				packet.setAddress(mcIPAddress);
				packet.setPort(multicastGroup.getMulticastPort());
				udpSocket.send(packet);
				udpSocket.close();
				
			}catch(Exception e){
				logger.error(e.getMessage());
			}
		}else{
			//Tcp not yet
		}
		
		
	}

}
