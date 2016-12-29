package com.lovi.quebic.sockm.messenger.impl;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.SocketTimeoutException;
import java.util.ArrayList;
import java.util.List;
import java.util.Random;
import java.util.UUID;
import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.async.AsyncExecutor;
import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.sockm.config.MulticastNetworkGroup;
import com.lovi.quebic.sockm.exception.SendMessageTimeOut;
import com.lovi.quebic.sockm.exception.SubscriberAddressNotFound;
import com.lovi.quebic.sockm.launcher.SockMLauncher;
import com.lovi.quebic.sockm.log.error.ErrorMessage;
import com.lovi.quebic.sockm.log.info.InfoMessage;
import com.lovi.quebic.sockm.message.messenger.Message;
import com.lovi.quebic.sockm.message.messenger.MessagePublish;
import com.lovi.quebic.sockm.message.messenger.MessageSend;
import com.lovi.quebic.sockm.message.messenger.MessageSendReply;
import com.lovi.quebic.sockm.messenger.Messenger;
import com.lovi.quebic.sockm.messenger.Publisher;
import com.lovi.quebic.sockm.messenger.Subscriber;
import com.lovi.quebic.sockm.messenger.SubscriberDirectory;

public class MessengerImpl implements Messenger{

	private final static Logger logger = LoggerFactory.getLogger(Messenger.class);
	
	//unique id for the Messanger
	private final String id;
	
	//SockMLauncher
	private SockMLauncher launcher;
	
	//Messanger
	private static Messenger instance;
	//Messanger lock
	private static Object lock = new Object();
	
	private Publisher publisher;
	
	private List<SubscriberDirectory> messageListeners;
	
	private MulticastNetworkGroup multicastGroup;
	
	private MessengerImpl(SockMLauncher launcher){
		this.id = genarateId();
		this.launcher = launcher;
		this.multicastGroup = launcher.getMulticastGroup();
		this.publisher = new PublisherImpl();
		this.messageListeners = launcher.getList("MESSAGE_LISTENER", SubscriberDirectory.class);
	}
	

	@Override
	public <I,O> Subscriber<I,O> subscribe(String address, Handler<Message<I,O>> handler) {
		Subscriber<I,O> subscriber = new SubscriberImpl<I,O>(publisher).subscribe(address, handler);
		logger.info(InfoMessage.SUBSCRIBER_START_LISTEN + address);
		
		if(launcher.getDiscoveryOption().isEnableMulticastMethod()){
			//Multicast mode no need to put TCP port
			messageListeners.add(new SubscriberDirectory(null, launcher.getLauncherId(), subscriber.getId(), address));
		}else{
			//Tcp discovery mode
		}
		return subscriber;
	}
	
	@Override
	public <I,O> Subscriber<I,O> subscribe(String address, Class<?> inputClassType, Class<?> outputClassType, Handler<Message<I,O>> handler) {
		Subscriber<I,O> subscriber = new SubscriberImpl<I,O>(publisher).subscribe(address, inputClassType, outputClassType, handler);
		logger.info(InfoMessage.SUBSCRIBER_START_LISTEN + address);
		
		if(launcher.getDiscoveryOption().isEnableMulticastMethod()){
			//Multicast mode no need to put TCP port
			messageListeners.add(new SubscriberDirectory(null, launcher.getLauncherId(), subscriber.getId(), address));
		}else{
			//Tcp discovery mode
		}
		return subscriber;
	}
	
	@Override
	public <I> void publish(String address, I message){
		
		AsyncExecutor<Boolean> asyncExecutor = AsyncExecutor.create();
		asyncExecutor.run(()->{
			
			publisher.publish(address, message);
			
			MessagePublish messagePublish = new MessagePublish();
			messagePublish.setSenderLauncherId_PUBLISH(launcher.getLauncherId());
			messagePublish.setAddress_PUBLISH(address);
			messagePublish.setValue_PUBLISH(message);
			
			publishRemote(address, messagePublish);
			
			return true;
			
		}, r->{
		}, f->{
		});
		
	}
	
	@Override
	public <I> void publishLocal(String address, I message){
		publisher.publish(address, message);
	}
	
	@Override
	public <I, O> void send(String address, I message, Handler<O> resultHandler, Handler<Throwable> failureHandler){
		send(address, message, null, null, resultHandler, failureHandler);
	}
	
	@Override
	public <I, O> void send(String address, I message, Class<?> inputClassType, Class<?> outputClassType,Handler<O> resultHandler, Handler<Throwable> failureHandler){
		
		AsyncExecutor<Message<I,O>> asyncExecutor = AsyncExecutor.create();
		asyncExecutor.run(()->{
			
			Message<I,O> msg = new Message<>(message);
			boolean isFoundSubscriber = publisher.send(address, msg);
			
			if(!isFoundSubscriber){
				//go for remote sending process
				sendRemote(address, msg);
				
			}
			return msg;
			
		}, msg->{
			
			if(!msg.isReplyFail()){
				if(resultHandler != null){
					
					try{
						
						//do test and convert
						if(outputClassType != null){
							O outputMessage = msg.getReply();
							
							ObjectMapper objectMapper = new ObjectMapper();
							
							if(!outputMessage.getClass().getName().equals(outputClassType.getName())){
								
								String jsonStr = objectMapper.writeValueAsString(outputMessage);
								outputMessage = (O) objectMapper.readValue(jsonStr, outputClassType);
								msg.reply(outputMessage);
							}
							
						}
						
						resultHandler.handle(msg.getReply());
					
					}catch(Exception e){
						
						if(failureHandler != null)
							failureHandler.handle(e);
						
					}
				}
			}else{
				if(failureHandler != null)
					failureHandler.handle(msg.getReplyFailure());
			}
			
		}, f->{
			
			if(failureHandler != null)
				failureHandler.handle(f);
			
		});
		
		
	}
	
	@Override
	public <I, O> boolean sendLocal(String subscriberId, Message<I,O> msg){
		return publisher.sendBySubscriberId(subscriberId, msg);
	}
	
	public static Messenger create(SockMLauncher launcher){
		
		synchronized (lock) {
			if(instance == null)
				instance = new MessengerImpl(launcher);
		}
		
		return instance;
	}
	
	private <T> void publishRemote(String address, T message){

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
	
	@SuppressWarnings("unchecked")
	private <I,O> void sendRemote(String address, Message<I,O> msg){

		if(launcher.getDiscoveryOption().isEnableMulticastMethod()){
			try{
				
				//Selecting subscriber id
				String subscriberId = selectsubscriberIdToSend(address);
				
				if(subscriberId == null)
					throw new SubscriberAddressNotFound(ErrorMessage.MESSENGER_SUBSCRIBER_ADDRESS_NOT_FOUND + address);
				
				MessageSend messageSend = new MessageSend();
				messageSend.setSenderLauncherId_SEND(launcher.getLauncherId());
				messageSend.setSubscriberId_SEND(subscriberId);
				messageSend.setValue_SEND(msg.getMessage());
				
				DatagramSocket udpSocket = new DatagramSocket();
		
				InetAddress mcIPAddress = InetAddress.getByName(multicastGroup.getMulticastAddress());
		
				ObjectMapper mapper = new ObjectMapper();
				byte[] messageBytes = mapper.writeValueAsString(messageSend).getBytes();
		
				DatagramPacket packet = new DatagramPacket(
						messageBytes
						,messageBytes.length);
				
				packet.setAddress(mcIPAddress);
				packet.setPort(multicastGroup.getMulticastPort());
				udpSocket.send(packet);
				
				DatagramPacket replyPacket = new DatagramPacket(new byte[9000],
						9000);
				replyPacket.setAddress(mcIPAddress);
				replyPacket.setPort(multicastGroup.getMulticastPort());
				
				try{
					//udpSocket.setSoTimeout(2000);
					udpSocket.receive(replyPacket);
					
					String replyMsgStr = new String(replyPacket.getData(),
							replyPacket.getOffset(), replyPacket.getLength());
					
					MessageSendReply messageSendReply = mapper.readValue(replyMsgStr, MessageSendReply.class);
					
					if(messageSendReply.getResponseThrowableMessage() == null)
						msg.reply((O) messageSendReply.getResponseValue());
					else
						msg.replyFailure(new Throwable(messageSendReply.getResponseThrowableMessage()));
					
					
				}catch(SocketTimeoutException e){
					msg.replyFailure(new SendMessageTimeOut(ErrorMessage.MESSENGER_SEND_MESSAGE_TIME_OUT + address));
				}catch (Exception e) {
					msg.replyFailure(e);
				}
				
				udpSocket.close();
				
			}catch(Exception e){
				msg.replyFailure(e);
			}
		}else{
			//Tcp Discovery mode not yet
		}
		
	}
	
	private String selectsubscriberIdToSend(String address){

		List<SubscriberDirectory> remoteMessageListener = new ArrayList<>();
		
		for(SubscriberDirectory messageListener : messageListeners){
			
			if(!launcher.getLauncherId().equals( messageListener.getLauncherId()) &&
					messageListener.getListenerAddress().equals(address)){
				remoteMessageListener.add(messageListener);
			}
			
		}
		
		if(remoteMessageListener.size() > 0){
			int randomIndex = new Random().nextInt(remoteMessageListener.size());
			return remoteMessageListener.get(randomIndex).getListenerId();
		}else
			return null;
	}
	
	private String genarateId(){
		return UUID.randomUUID().toString();
	}
	
	@Override
	public String getLauncherId(){
		return this.id;
	}
}
