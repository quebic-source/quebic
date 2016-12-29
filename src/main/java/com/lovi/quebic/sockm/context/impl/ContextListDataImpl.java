package com.lovi.quebic.sockm.context.impl;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.sockm.config.MulticastNetworkGroup;
import com.lovi.quebic.sockm.context.ContextListData;
import com.lovi.quebic.sockm.launcher.SockMLauncher;
import com.lovi.quebic.sockm.log.info.InfoMessage;
import com.lovi.quebic.sockm.message.list.MessageGetList;
import com.lovi.quebic.sockm.util.SharedList;

public class ContextListDataImpl implements ContextListData {

	private final Logger logger = LoggerFactory.getLogger(ContextListData.class);
	
	private SockMLauncher launcher;
	private MulticastNetworkGroup multicastGroup;
	private Map<String, List<?>> dataList = new HashMap<>();
	
	public ContextListDataImpl(SockMLauncher launcher) {
		this.launcher = launcher;
		multicastGroup = launcher.getMulticastGroup();
	}
	
	@Override
	public <E> List<E> getList(String key){
		
		if(dataList.containsKey(key)){
			@SuppressWarnings("unchecked")
			List<E> list = (List<E>) dataList.get(key);
			return list;
		}
		else{
			List<E> list = createtListData(key, null);
			dataList.put(key, list);
			return list;
		}
		
	}
	
	@Override
	public <E> List<E> getList(String key, Class<E> typeOfValue){
		
		if(dataList.containsKey(key)){
			@SuppressWarnings("unchecked")
			List<E> list = (List<E>) dataList.get(key);
			return list;
		}
		else{
			List<E> list = createtListData(key, typeOfValue);
			dataList.put(key, list);
			return list;
		}
		
	}
	
	@Override
	public void updateListForAdd(String contextKey, Object value){
		if(dataList.containsKey(contextKey)){
			@SuppressWarnings("unchecked")
			SharedList<Object> list = (SharedList<Object>) dataList.get(contextKey);
			list._add(value);
		}
	}
	
	@Override
	public void updateListForAdd(String contextKey, int index, Object value){
		if(dataList.containsKey(contextKey)){
			@SuppressWarnings("unchecked")
			SharedList<Object> list = (SharedList<Object>) dataList.get(contextKey);
			list._add(index, value);
		}
	}
	
	@Override
	public void updateListForAddAll(String contextKey, Collection<?> collection){
		if(dataList.containsKey(contextKey)){
			@SuppressWarnings("unchecked")
			SharedList<Object> list = (SharedList<Object>) dataList.get(contextKey);
			list._addAll(collection);
		}
	}
	
	@Override
	public void updateListForAddAll(String contextKey, int index ,Collection<?> collection){
		if(dataList.containsKey(contextKey)){
			@SuppressWarnings("unchecked")
			SharedList<Object> list = (SharedList<Object>) dataList.get(contextKey);
			list._addAll(index, collection);
		}
	}
	
	@Override
	public void updateListForRemove(String contextKey, Object o){
		if(dataList.containsKey(contextKey)){
			@SuppressWarnings("unchecked")
			SharedList<Object> list = (SharedList<Object>) dataList.get(contextKey);
			list._remove(o);
		}
	}
	
	@Override
	public void updateListForRemoveAll(String contextKey, Collection<?> collection){
		if(dataList.containsKey(contextKey)){
			@SuppressWarnings("unchecked")
			SharedList<Object> list = (SharedList<Object>) dataList.get(contextKey);
			list._removeAll(collection);
		}
	}
	
	@Override
	public void updateListForClear(String contextKey){
		if(dataList.containsKey(contextKey)){
			@SuppressWarnings("unchecked")
			SharedList<Object> list = (SharedList<Object>) dataList.get(contextKey);
			list._clear();
		}
	}
	
	@SuppressWarnings("unchecked")
	private <E> List<E> createtListData(String contextKey, Class<E> typeOfValue){
		SharedList<E> list = new SharedList<>(launcher, contextKey, typeOfValue);
		
		try {
			
			DatagramSocket udpSocket = new DatagramSocket();
			
			InetAddress mcIPAddress = InetAddress.getByName(multicastGroup.getMulticastAddress());
			
			MessageGetList messageGetList = new MessageGetList();
			messageGetList.setSenderLauncherId_GET_LIST(launcher.getLauncherId());
			messageGetList.setContextKey_GET_LIST(contextKey);

			ObjectMapper mapper = new ObjectMapper();
			byte[] joinRequestMessageBytes = mapper.writeValueAsString(
					messageGetList).getBytes();
			
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

				messageGetList = mapper.readValue(replyMsgStr, MessageGetList.class);
				
				list._addAll((Collection<? extends E>) messageGetList.getContextDataList_GET_LIST());

			} catch (Exception e) {
				logger.info(InfoMessage.CREATE_NEW_LIST);
				//logger.error(e.getMessage());
			}
			
			udpSocket.close();
			
		} catch (Exception e) {
			logger.error("createtMapData " + e.getMessage());
		}
		return list;
	}
	
	@Override
	public Map<String, List<?>> getDataList(){
		return dataList;
	}

}
