package com.lovi.quebic.sockm.server;

import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.MulticastSocket;
import java.util.Collection;
import java.util.List;
import java.util.Map;

import org.codehaus.jackson.JsonGenerationException;
import org.codehaus.jackson.map.JsonMappingException;
import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.async.AsyncExecutor;
import com.lovi.quebic.sockm.config.MulticastNetworkGroup;
import com.lovi.quebic.sockm.context.ContextListData;
import com.lovi.quebic.sockm.context.ContextMapData;
import com.lovi.quebic.sockm.exception.MulticastServerException;
import com.lovi.quebic.sockm.launcher.SockMLauncher;
import com.lovi.quebic.sockm.message.list.MessageAddAllList;
import com.lovi.quebic.sockm.message.list.MessageAddList;
import com.lovi.quebic.sockm.message.list.MessageClearList;
import com.lovi.quebic.sockm.message.list.MessageGetList;
import com.lovi.quebic.sockm.message.list.MessageRemoveAllList;
import com.lovi.quebic.sockm.message.list.MessageRemoveList;
import com.lovi.quebic.sockm.message.map.MessageClearMap;
import com.lovi.quebic.sockm.message.map.MessageGetMap;
import com.lovi.quebic.sockm.message.map.MessagePutAllMap;
import com.lovi.quebic.sockm.message.map.MessagePutMap;
import com.lovi.quebic.sockm.message.map.MessageRemoveMap;
import com.lovi.quebic.sockm.message.messenger.Message;
import com.lovi.quebic.sockm.message.messenger.MessagePublish;
import com.lovi.quebic.sockm.message.messenger.MessageSend;
import com.lovi.quebic.sockm.message.messenger.MessageSendReply;
import com.lovi.quebic.sockm.messenger.Messenger;

public class MulticastServer {

	private final Logger logger = LoggerFactory
			.getLogger(MulticastServer.class);

	private MulticastNetworkGroup multicastGroup;
	private SockMLauncher launcher;
	
	private Thread serverThread;
	
	private ContextMapData contextMapData;
	private ContextListData contextListData;
	
	private Messenger messenger;

	public MulticastServer(SockMLauncher launcher, ContextMapData contextMapData, ContextListData contextListData, Messenger messenger) {
		this.launcher = launcher;
		this.multicastGroup = this.launcher.getMulticastGroup();
		this.contextMapData = contextMapData;
		this.contextListData = contextListData;
		this.messenger = messenger;
	}

	public void listen() throws MulticastServerException {

		final MulticastSocket mcSocket;
		final InetAddress mcIPAddress;

		try {
			mcIPAddress = InetAddress.getByName(multicastGroup
					.getMulticastAddress());
			mcSocket = new MulticastSocket(multicastGroup.getMulticastPort());
			mcSocket.joinGroup(mcIPAddress);


		} catch (Exception e) {
			throw new MulticastServerException(e.getMessage());
		}

		serverThread = new Thread(new Runnable() {

			@Override
			public void run() {
				try {

					while (true) {
						
						//check for stop
						if(Thread.interrupted())
							break;
						
						DatagramPacket packet = new DatagramPacket(new byte[1024], 1024);
						mcSocket.receive(packet);
						
						AsyncExecutor<Boolean> asyncExecutor = AsyncExecutor.create();
						asyncExecutor.run(()->{
							
							try{
								
								String msg = new String(packet.getData(),
										packet.getOffset(), packet.getLength());

								ObjectMapper objectMapper = new ObjectMapper();
								DatagramSocket udpSocket = new DatagramSocket();
								
								/**
								 * Map Processing
								 */
								// MessageGetMap
								try (MessageGetMap message = objectMapper.readValue(
										msg, MessageGetMap.class)) {
									processMessageGetMap(objectMapper, message,
											udpSocket, packet, mcIPAddress);
									return true;
								} catch (Exception e) {
								}

								// MessagePutMap
								try (MessagePutMap message = objectMapper.readValue(
										msg, MessagePutMap.class)) {
									processMessagePutMap(message);
									return true;
								} catch (Exception e) {
								}

								// MessagePutAllMap
								try (MessagePutAllMap message = objectMapper.readValue(
										msg, MessagePutAllMap.class)) {
									processMessagePutAllMap(message);
									return true;
								} catch (Exception e) {
								}

								// MessageRemoveMap
								try (MessageRemoveMap message = objectMapper.readValue(
										msg, MessageRemoveMap.class)) {
									processMessageRemoveMap(message);
									return true;
								} catch (Exception e) {
								}

								// MessageClearMap
								try (MessageClearMap message = objectMapper.readValue(
										msg, MessageClearMap.class)) {
									processMessageClearMap(message);
									return true;
								} catch (Exception e) {
								}
								/**
								 * End Map Processing
								 */

								/**
								 * List Processing
								 */
								// MessageGetMap
								try (MessageGetList message = objectMapper.readValue(
										msg, MessageGetList.class)) {
									processMessageGetList(objectMapper, message,
											udpSocket, packet, mcIPAddress);
									return true;
								} catch (Exception e) {
								}
								
								//MessageAddList
								try (MessageAddList message = objectMapper.readValue(
										msg, MessageAddList.class)) {
									processMessageAddList(message);
									return true;
								} catch (Exception e) {
								}
								
								//MessageAddAllList
								try (MessageAddAllList message = objectMapper.readValue(
										msg, MessageAddAllList.class)) {
									processMessageAddAllList(message);
									return true;
								} catch (Exception e) {
								}
								
								//MessageRemoveList
								try (MessageRemoveList message = objectMapper.readValue(
										msg, MessageRemoveList.class)) {
									processMessageRemoveList(message);
									return true;
								} catch (Exception e) {
								}
								
								//MessageRemoveAllList
								try (MessageRemoveAllList message = objectMapper.readValue(
										msg, MessageRemoveAllList.class)) {
									processMessageRemoveAllList(message);
									return true;
								} catch (Exception e) {
								}
								
								//MessageClearList
								try (MessageClearList message = objectMapper.readValue(
										msg, MessageClearList.class)) {
									processMessageClearList(message);
									return true;
								} catch (Exception e) {
								}
								
								/**
								 * End List Processing
								 */
								
								/**
								 * Start Messenger Publish Processing
								 */
								//MessageClearList
								try (MessagePublish message = objectMapper.readValue(
										msg, MessagePublish.class)) {
									processMessagePublish(message);
									return true;
								} catch (Exception e) {
								}
								
								/**
								 * End Messenger Publish Processing
								 */
								
								/**
								 * Start Messenger Send Processing
								 */
								//MessageClearList
								try (MessageSend message = objectMapper.readValue(
										msg, MessageSend.class)) {
									processMessageSend(objectMapper,
											udpSocket, packet, mcIPAddress,message);
									return true;
								} catch (Exception e) {
								}
								
								/**
								 * End Messenger Send Processing
								 */
								
							}catch(Exception e){
							}
							
							return true;
							
						}, r->{
						}, f->{	
						});
						
					}

					mcSocket.leaveGroup(mcIPAddress);
					mcSocket.close();

				} catch (Exception e) {
					logger.error("multiCastListen " + e.getMessage());
				}

			}

		}, "sockm-server-thread");

		serverThread.start();

	}

	/*
	 * Map Processing
	 */
	
	private void processMessageGetMap(ObjectMapper objectMapper,
			MessageGetMap message, DatagramSocket udpSocket,
			DatagramPacket packet, InetAddress mcIPAddress)
			throws JsonGenerationException, JsonMappingException, IOException {

		// check this is not me
		if (!message.getSenderLauncherId_GET_MAP().equals(
				launcher.getLauncherId())) {

			InetAddress receivedAddressUDP = packet.getAddress();
			int receivedPortUDP = packet.getPort();

			String contextKey = message.getContextKey_GET_MAP();

			MessageGetMap replyMessage = new MessageGetMap();
			replyMessage.setSenderLauncherId_GET_MAP(launcher.getLauncherId());
			replyMessage.setContextKey_GET_MAP(contextKey);

			Map<?, ?> map = contextMapData.getDataMap()
					.get(contextKey);
			if (map != null) {
				replyMessage.setContextDataMap_GET_MAP(map);
				
				System.out.println(objectMapper.writeValueAsString(map));
				
				byte[] replyMegBytes = objectMapper.writeValueAsString(
						replyMessage).getBytes();

				DatagramPacket replyPacket = new DatagramPacket(replyMegBytes,
						replyMegBytes.length);
				replyPacket.setAddress(receivedAddressUDP);
				replyPacket.setPort(receivedPortUDP);
				udpSocket.send(replyPacket);

			}

		}

	}

	private void processMessagePutMap(MessagePutMap message) {
		// check this is not me
		if (!message.getSenderLauncherId_PUT_MAP().equals(
				launcher.getLauncherId())) {

			String contextKey = message.getContextKey_PUT_MAP();
			Object key = message.getMapKey_PUT_MAP();
			Object value = message.getMapValue_PUT_MAP();
			contextMapData.updateMapForPut(contextKey, key, value);

		}
	}

	private void processMessagePutAllMap(MessagePutAllMap message) {
		// check this is not me
		if (!message.getSenderLauncherId_PUTALL_MAP().equals(
				launcher.getLauncherId())) {

			String contextKey = message.getContextKey_PUTALL_MAP();
			Map<?, ?> map = message.getMap_PUTALL_MAP();
			contextMapData.updateMapForPutAll(contextKey, map);

		}
	}

	private void processMessageRemoveMap(MessageRemoveMap message) {
		// check this is not me
		if (!message.getSenderLauncherId_REMOVE_MAP().equals(
				launcher.getLauncherId())) {

			String contextKey = message.getContextKey_REMOVE_MAP();
			Object key = message.getMapKey_REMOVE_MAP();
			contextMapData.updateMapForRemove(contextKey, key);

		}
	}

	private void processMessageClearMap(MessageClearMap message) {
		// check this is not me
		if (!message.getSenderLauncherId_CLEAR_MAP().equals(
				launcher.getLauncherId())) {

			String contextKey = message.getContextKey_CLEAR_MAP();
			contextMapData.updateMapForClear(contextKey);

		}
	}
	
	/*
	 * End Map Processing
	 */
	

	/*
	 * List Processing
	 */
	private void processMessageGetList(ObjectMapper objectMapper,
			MessageGetList message, DatagramSocket udpSocket,
			DatagramPacket packet, InetAddress mcIPAddress) throws JsonGenerationException, JsonMappingException, IOException {

		// check this is not me
		if (!message.getSenderLauncherId_GET_LIST().equals(
				launcher.getLauncherId())) {

			InetAddress receivedAddressUDP = packet.getAddress();
			int receivedPortUDP = packet.getPort();

			String contextKey = message.getContextKey_GET_LIST();

			MessageGetList replyMessage = new MessageGetList();
			replyMessage.setSenderLauncherId_GET_LIST(launcher.getLauncherId());
			replyMessage.setContextKey_GET_LIST(contextKey);

			List<?> list = contextListData.getDataList().get(contextKey);
			if (list != null) {
				replyMessage.setContextDataList_GET_LIST(list);

				byte[] replyMegBytes = objectMapper.writeValueAsString(
						replyMessage).getBytes();

				DatagramPacket replyPacket = new DatagramPacket(replyMegBytes,
						replyMegBytes.length);
				replyPacket.setAddress(receivedAddressUDP);
				replyPacket.setPort(receivedPortUDP);
				udpSocket.send(replyPacket);

			}

		}

	}
	
	private void processMessageAddList(MessageAddList message) {
		// check this is not me
		if (!message.getSenderLauncherId_ADD_LIST().equals(
				launcher.getLauncherId())) {

			String contextKey = message.getContextKey_ADD_LIST();
			int index = message.getListIndex_ADD_LIST();
			Object value = message.getListValue_ADD_LIST();
			
			if(index == -1)
				contextListData.updateListForAdd(contextKey, value);
			else
				contextListData.updateListForAdd(contextKey, index, value);

		}
	}
	
	private void processMessageAddAllList(MessageAddAllList message) {
		// check this is not me
		if (!message.getSenderLauncherId_ADDALL_LIST().equals(
				launcher.getLauncherId())) {

			String contextKey = message.getContextKey_ADDALL_LIST();
			int index = message.getListIndex_ADDALL_LIST();
			Collection<?> collection = message.getCollection_ADDALL_LIST();
			
			if(index == -1)
				contextListData.updateListForAddAll(contextKey, collection);
			else
				contextListData.updateListForAddAll(contextKey, index, collection);

		}
	}
	
	private void processMessageRemoveList(MessageRemoveList message) {
		// check this is not me
		if (!message.getSenderLauncherId_REMOVE_LIST().equals(
				launcher.getLauncherId())) {

			String contextKey = message.getContextKey_REMOVE_LIST();
			Object o = message.getListValue_REMOVE_LIST();
			contextListData.updateListForRemove(contextKey, o);

		}
	}
	
	private void processMessageRemoveAllList(MessageRemoveAllList message) {
		// check this is not me
		if (!message.getSenderLauncherId_REMOVEALL_LIST().equals(
				launcher.getLauncherId())) {

			String contextKey = message.getContextKey_REMOVEALL_LIST();
			Collection<?> collection = message.getCollection_REMOVEALL_LIST();
			contextListData.updateListForRemoveAll(contextKey, collection);

		}
	}

	private void processMessageClearList(MessageClearList message) {
		// check this is not me
		if (!message.getSenderLauncherId_CLEAR_LIST().equals(
				launcher.getLauncherId())) {

			String contextKey = message.getContextKey_CLEAR_LIST();
			contextListData.updateListForClear(contextKey);

		}
	}
	/*
	 * End List Processing
	 */
	
	
	/*
	 * Message Publish Processing
	 */
	private void processMessagePublish(MessagePublish message) {
		
		if(!message.getSenderLauncherId_PUBLISH().equals(launcher.getLauncherId())){
			AsyncExecutor<Boolean> asyncExecutor = AsyncExecutor.create();
			asyncExecutor.run(()->{
				messenger.publishLocal(message.getAddress_PUBLISH(), message.getValue_PUBLISH());
				return true;
			}, r->{
			}, f->{
			});
		}
	}
	/*
	 * End Message Publish Processing
	 */
	

	/*
	 * Message Send Processing
	 */
	private void processMessageSend(ObjectMapper objectMapper, DatagramSocket udpSocket,
			DatagramPacket packet, InetAddress mcIPAddress, MessageSend message) throws JsonGenerationException, JsonMappingException, IOException {
		
		if(!message.getSenderLauncherId_SEND().equals(launcher.getLauncherId())){
			
			Message<?,?> msg = new Message<>(message.getValue_SEND());
			
			if(messenger.sendLocal(message.getSubscriberId_SEND(), msg)){
			
				InetAddress receivedAddressUDP = packet.getAddress();
				int receivedPortUDP = packet.getPort();
				
				MessageSendReply messageSendReply = new MessageSendReply();
				
				if(!msg.isReplyFail())
					messageSendReply.setResponseValue(msg.getReply());
				else
					messageSendReply.setResponseThrowableMessage(msg.getReplyFailure().getMessage());
				
				byte[] replyMegBytes = objectMapper.writeValueAsString(
						messageSendReply).getBytes();
	
				DatagramPacket replyPacket = new DatagramPacket(replyMegBytes,
						replyMegBytes.length);
				replyPacket.setAddress(receivedAddressUDP);
				replyPacket.setPort(receivedPortUDP);
				udpSocket.send(replyPacket);
				
			}
		}
	}
	/*
	 * End Message Send Processing
	 */
	
	public synchronized void stopServer() {
		if(serverThread != null)
			serverThread.interrupt();
	}
}
