package com.lovi.quebic.sockm.util;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Iterator;
import java.util.List;
import java.util.ListIterator;

import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.sockm.config.MulticastNetworkGroup;
import com.lovi.quebic.sockm.launcher.SockMLauncher;
import com.lovi.quebic.sockm.message.list.MessageAddAllList;
import com.lovi.quebic.sockm.message.list.MessageAddList;
import com.lovi.quebic.sockm.message.list.MessageClearList;
import com.lovi.quebic.sockm.message.list.MessageRemoveAllList;
import com.lovi.quebic.sockm.message.list.MessageRemoveList;

public class SharedList<E> implements List<E>{

	private final Logger logger = LoggerFactory.getLogger(SharedList.class);
	
	private Class<E> typeOfValue;
	
	private List<E> list = new ArrayList<>();
	
	private SockMLauncher launcher;
	private String contextKey;
	private MulticastNetworkGroup multicastGroup;
	
	public SharedList(SockMLauncher launcher, String contextKey, Class<E> typeOfValue) {
		this.launcher = launcher;
		this.multicastGroup = launcher.getMulticastGroup();
		this.contextKey = contextKey;
		this.typeOfValue = typeOfValue;
	}
	
	@Override
	public int size() {
		return list.size();
	}

	@Override
	public boolean isEmpty() {
		return list.isEmpty();
	}

	@Override
	public boolean contains(Object o) {
		return list.contains(o);
	}

	@Override
	public Iterator<E> iterator() {
		return list.iterator();
	}

	@Override
	public Object[] toArray() {
		return list.toArray();
	}

	@Override
	public <T> T[] toArray(T[] a) {
		return list.toArray(a);
	}

	@Override
	public boolean add(E e) {
		
		MessageAddList addList = new MessageAddList();
		addList.setSenderLauncherId_ADD_LIST(launcher.getLauncherId());
		addList.setContextKey_ADD_LIST(contextKey);
		addList.setListValue_ADD_LIST(e);
		
		sendMessage(addList);
		
		return list.add(e);
	}
	
	public boolean _add(E e) {
		
		ObjectMapper objectMapper = new ObjectMapper();
		
		if(typeOfValue != null){
			try {
				
				if(!e.getClass().getName().equals(typeOfValue.getName())){
					String jsonStr = objectMapper.writeValueAsString(e);
					e = objectMapper.readValue(jsonStr, typeOfValue);
				}
				
			} catch (Exception ex) {
				logger.error("_add process value " + ex.getMessage());
			}

		}
		
		return list.add(e);
	}

	@Override
	public boolean remove(Object o) {
		
		MessageRemoveList removeList = new MessageRemoveList();
		removeList.setSenderLauncherId_REMOVE_LIST(launcher.getLauncherId());
		removeList.setContextKey_REMOVE_LIST(contextKey);
		removeList.setListValue_REMOVE_LIST(o);
		
		sendMessage(removeList);
		
		return list.remove(o);
	}
	
	public boolean _remove(Object o) {
		
		ObjectMapper objectMapper = new ObjectMapper();
		
		if(typeOfValue != null){
			try {
				
				if(!o.getClass().getName().equals(typeOfValue.getName())){
					String jsonStr = objectMapper.writeValueAsString(o);
					o = objectMapper.readValue(jsonStr, typeOfValue);
				}
				
			} catch (Exception ex) {
				logger.error("_remove process value " + ex.getMessage());
			}

		}
		
		
		return list.remove(o);
	}

	@Override
	public boolean containsAll(Collection<?> c) {
		return list.containsAll(c);
	}

	@Override
	public boolean addAll(Collection<? extends E> c) {
		
		MessageAddAllList addAllList = new MessageAddAllList();
		addAllList.setSenderLauncherId_ADDALL_LIST(launcher.getLauncherId());
		addAllList.setContextKey_ADDALL_LIST(contextKey);
		addAllList.setCollection_ADDALL_LIST(c);
		
		sendMessage(addAllList);
		
		return list.addAll(c);
	}
	
	public boolean _addAll(Collection<? extends E> c) {
		
		ObjectMapper objectMapper = new ObjectMapper();
		Collection<E> tmpList = new ArrayList<>();
		if(typeOfValue != null){

			for(E e : c){
				try {
					
					if(!e.getClass().getName().equals(typeOfValue.getName())){
						String jsonStr = objectMapper.writeValueAsString(e);
						e = objectMapper.readValue(jsonStr, typeOfValue);
					}
					
				} catch (Exception ex) {
					logger.error("_addAll process value " + ex.getMessage());
				}
				
				tmpList.add(e);
				
			}
			return list.addAll(tmpList);
		}else{

			return list.addAll(c);
		}
		
	}

	@Override
	public boolean addAll(int index, Collection<? extends E> c) {
		
		MessageAddAllList addAllList = new MessageAddAllList();
		addAllList.setSenderLauncherId_ADDALL_LIST(launcher.getLauncherId());
		addAllList.setContextKey_ADDALL_LIST(contextKey);
		addAllList.setCollection_ADDALL_LIST(c);
		addAllList.setListIndex_ADDALL_LIST(index);
		
		sendMessage(addAllList);
		
		return list.addAll(index, c);
	}
	
	public boolean _addAll(int index, Collection<? extends E> c) {
		
		ObjectMapper objectMapper = new ObjectMapper();
		Collection<E> tmpList = new ArrayList<>();
		if(typeOfValue != null){

			for(E e : c){
				try {
					
					if(!e.getClass().getName().equals(typeOfValue.getName())){
						String jsonStr = objectMapper.writeValueAsString(e);
						e = objectMapper.readValue(jsonStr, typeOfValue);
					}
					
				} catch (Exception ex) {
					logger.error("_addAll process value " + ex.getMessage());
				}
				
				tmpList.add(e);
				
			}
			return list.addAll(index, tmpList);
		}else{

			return list.addAll(index, c);
		}
	}

	@Override
	public boolean removeAll(Collection<?> c) {
		
		MessageRemoveAllList removeAllList = new MessageRemoveAllList();
		removeAllList.setSenderLauncherId_REMOVEALL_LIST(launcher.getLauncherId());
		removeAllList.setContextKey_REMOVEALL_LIST(contextKey);
		removeAllList.setCollection_REMOVEALL_LIST(c);
		
		sendMessage(removeAllList);
		
		return list.removeAll(c);
	}
	
	public boolean _removeAll(Collection<?> c) {
		
		ObjectMapper objectMapper = new ObjectMapper();
		Collection<Object> tmpList = new ArrayList<>();
		if(typeOfValue != null){

			for(Object e : c){
				try {
					
					if(!e.getClass().getName().equals(typeOfValue.getName())){
						String jsonStr = objectMapper.writeValueAsString(e);
						e = objectMapper.readValue(jsonStr, typeOfValue);
					}
					
				} catch (Exception ex) {
					logger.error("_addAll process value " + ex.getMessage());
				}
				
				tmpList.add(e);
				
			}
			return list.removeAll(tmpList);
		}else{
			return list.removeAll(c);
		}
		
	}

	@Override
	public boolean retainAll(Collection<?> c) {
		return list.retainAll(c);
	}

	@Override
	public void clear() {
		
		MessageClearList clearList = new MessageClearList();
		clearList.setSenderLauncherId_CLEAR_LIST(launcher.getLauncherId());
		clearList.setContextKey_CLEAR_LIST(contextKey);
		
		sendMessage(clearList);
		
		list.clear();
	}
	
	public void _clear() {
		list.clear();
	}

	@Override
	public E get(int index) {
		return list.get(index);
	}

	@Override
	public E set(int index, E element) {
		return list.set(index, element);
	}

	@Override
	public void add(int index, E element) {
		
		MessageAddList addList = new MessageAddList();
		addList.setSenderLauncherId_ADD_LIST(launcher.getLauncherId());
		addList.setContextKey_ADD_LIST(contextKey);
		addList.setListValue_ADD_LIST(element);
		addList.setListIndex_ADD_LIST(index);
		
		sendMessage(addList);
		
		list.add(index, element);
	}
	
	public void _add(int index, E element) {
		
		ObjectMapper objectMapper = new ObjectMapper();
		
		if(typeOfValue != null){
			try {
				
				if(!element.getClass().getName().equals(typeOfValue.getName())){
					String jsonStr = objectMapper.writeValueAsString(element);
					element = objectMapper.readValue(jsonStr, typeOfValue);
				}
				
			} catch (Exception ex) {
				logger.error("_add process value " + ex.getMessage());
			}

		}
		
		list.add(index, element);
	}

	@Override
	public E remove(int index) {
		E o = get(index);
		remove(o);
		return o;
	}

	@Override
	public int indexOf(Object o) {
		return list.indexOf(o);
	}

	@Override
	public int lastIndexOf(Object o) {
		return list.lastIndexOf(o);
	}

	@Override
	public ListIterator<E> listIterator() {
		return list.listIterator();
	}

	@Override
	public ListIterator<E> listIterator(int index) {
		return list.listIterator(index);
	}

	@Override
	public List<E> subList(int fromIndex, int toIndex) {
		return list.subList(fromIndex, toIndex);
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

	@Override
	public String toString() {
		return "SharedList [list=" + list + ", contextKey=" + contextKey + "]";
	}

}
