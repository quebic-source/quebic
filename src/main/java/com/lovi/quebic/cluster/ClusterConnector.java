package com.lovi.quebic.cluster;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.MulticastSocket;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.async.Future;
import com.lovi.quebic.cluster.loadblancer.LoadBlancer;
import com.lovi.quebic.cluster.message.CoreMessage;
import com.lovi.quebic.cluster.message.MessageKey;
import com.lovi.quebic.cluster.option.MulticastGroup;
import com.lovi.quebic.common.Message;
import com.lovi.quebic.web.ApplicationContextData;
import com.lovi.quebic.web.SessionStore;

public class ClusterConnector {

	private final static Logger logger = LoggerFactory.getLogger(ClusterConnector.class);
	
	private List<Member> tcpMembers = new ArrayList<>();
	private Member localTcpMember;
	private LoadBlancer loadBlancer;
	
	private final String multicastAddress;
	private final int multicastPort;
	
	private ApplicationContextData applicationContextData;
	private SessionStore sessionStore;
	
	/**
	 * <b>default multicastAddress</b> = 230.1.1.1
	 * </br>
	 * <b>default multicastPort</b> = 12345
	 */
	public ClusterConnector() {
		multicastAddress = "230.1.1.1";
		multicastPort = 12345;
	}
	
	public ClusterConnector(String multicastAddress, int multicastPort) {
		this.multicastAddress = multicastAddress;
		this.multicastPort = multicastPort;
	}
	
	
	public void startAsMaster(String address, int port, LoadBlancer loadBlancer, Future<String> future) throws Exception{
		
		localTcpMember = new Member(address, port, true);
		this.loadBlancer = loadBlancer;
		this.loadBlancer.setMembers(tcpMembers);

		
		
		new Thread(()->{
			ProxyServer proxyServer = new ProxyServer(address, port, this.loadBlancer);
			try {
				proxyServer.start(future);

				if (future != null){
					if(future.isFail()){
						throw future.getFailure();
					}
				}
				
			} catch (Throwable e) {
				if (future != null)
					future.setFailure(e);
			}
		}).start();
		
		
		joinGroup();
		multiCastListen();
		UnreachableMembersChecker.startCheck(this);
		
	}
	
	public void startAsWorker(String workerAddress, int workerPort) {
		
		localTcpMember = new Member(workerAddress, workerPort);
		tcpMembers.add(localTcpMember);
		
		joinGroup();
		getLiveData();
		multiCastListen();
		UnreachableMembersChecker.startCheck(this);
	}
	
	private void joinGroup(){

		try {

			DatagramSocket udpSocket = new DatagramSocket();

			InetAddress mcIPAddress = InetAddress.getByName(multicastAddress);

			CoreMessage joinRequestMessage = new CoreMessage();
			joinRequestMessage.setKey(MessageKey.JOIN_NEW_SERVER_MEMBER);
			joinRequestMessage.setMember(localTcpMember);

			ObjectMapper mapper = new ObjectMapper();
			byte[] joinRequestMessageBytes = mapper.writeValueAsString(
					joinRequestMessage).getBytes();

			DatagramPacket joinRequestPacket = new DatagramPacket(
					joinRequestMessageBytes, joinRequestMessageBytes.length);
			joinRequestPacket.setAddress(mcIPAddress);
			joinRequestPacket.setPort(multicastPort);
			udpSocket.send(joinRequestPacket);

			DatagramPacket replyPacket = new DatagramPacket(new byte[1024],
					1024);
			replyPacket.setAddress(mcIPAddress);
			replyPacket.setPort(multicastPort);

			try {
				logger.info(Message.CLUSTER_MEMBER_WAITING_JOIN.getMessage());
				udpSocket.setSoTimeout(2000);
				udpSocket.receive(replyPacket);

				String replyMsgStr = new String(replyPacket.getData(),
						replyPacket.getOffset(), replyPacket.getLength());

				CoreMessage replyMessage = mapper.readValue(replyMsgStr,
						CoreMessage.class);

				if (replyMessage.getKey().equals(MessageKey.WELCOME_NEW_SERVER_MEMBER)) {
					try {
						tcpMembers.addAll(replyMessage.getMembers());
						
					} catch (Exception e) {
						System.err.println("unable to cast memebers list");
					}
				}

			} catch (Exception e) {
				logger.info(Message.CLUSTER_FIRST_MEMBER.getMessage());
			}
			printTCPMembers();
			udpSocket.close();
		} catch (Exception e) {
			logger.error("joinGroup " + e.getMessage());
		}

	}
	
	/**
	 * get data(Session, Application) from live clusters 
	 */
	private void getLiveData(){

		try {

			DatagramSocket udpSocket = new DatagramSocket();

			InetAddress mcIPAddress = InetAddress.getByName(multicastAddress);

			CoreMessage requestMessage = new CoreMessage();
			requestMessage.setKey(MessageKey.REQUEST_LIVE_DATA);

			ObjectMapper mapper = new ObjectMapper();
			byte[] requestMessageBytes = mapper.writeValueAsString(
					requestMessage).getBytes();

			DatagramPacket requestPacket = new DatagramPacket(
					requestMessageBytes, requestMessageBytes.length);
			requestPacket.setAddress(mcIPAddress);
			requestPacket.setPort(multicastPort);
			udpSocket.send(requestPacket);

			DatagramPacket replyPacket = new DatagramPacket(new byte[1024],
					1024);
			replyPacket.setAddress(mcIPAddress);
			replyPacket.setPort(multicastPort);

			try {
				
				logger.info(Message.CLUSTER_WORKER_MEMBER_WAITING_LIVE_DATA.getMessage());
				udpSocket.setSoTimeout(2000);
				udpSocket.receive(replyPacket);

				String replyMsgStr = new String(replyPacket.getData(),
						replyPacket.getOffset(), replyPacket.getLength());

				CoreMessage replyMessage = mapper.readValue(replyMsgStr,
						CoreMessage.class);

				if (replyMessage.getKey().equals(MessageKey.LIVE_DATA)) {
					try {
						sessionStore.setUsers(replyMessage.getSessionParm());
						applicationContextData.setParms(replyMessage.getApplicationContextParms());
					} catch (Exception e) {
						System.err.println("unable to get session store " + e.getMessage());
					}
				}

			} catch (Exception e) {
				logger.info(Message.CLUSTER_FIRST_WORKER_MEMBER.getMessage());
			}
	
			udpSocket.close();
		} catch (Exception e) {
			logger.error("getLiveData " + e.getMessage());
		}

	}
	
	private void multiCastListen() {
		
		Thread thread = new Thread(new Runnable() {
			
			@Override
			public void run() {
				try {

					InetAddress mcIPAddress = InetAddress.getByName(multicastAddress);
					MulticastSocket mcSocket = new MulticastSocket(multicastPort);
					mcSocket.joinGroup(mcIPAddress);
					
					DatagramPacket packet = new DatagramPacket(new byte[1024], 1024);

					while (true) {
						mcSocket.receive(packet);
						String msg = new String(packet.getData(), packet.getOffset(), packet.getLength());

						ObjectMapper objectMapper = new ObjectMapper();
						
						try(CoreMessage message = objectMapper.readValue(msg, CoreMessage.class)){
							
							if (message.getKey().equals(MessageKey.JOIN_NEW_SERVER_MEMBER)) {
								
								processNewMemberJoin(message, objectMapper, mcSocket, packet);
								
							}else if (message.getKey().equals(MessageKey.REQUEST_LIVE_DATA) && !localTcpMember.isMaster()) {
								
								processRequestLiveData(message, objectMapper, mcSocket, packet);
								
							}
							
							//application context data
							else if (message.getKey().equals(MessageKey.APPLICATION_PARM_PUT) && !localTcpMember.isMaster()) {
								
								processApplicationParmPut(message);
								
							}else if (message.getKey().equals(MessageKey.APPLICATION_PARM_REMOVE) && !localTcpMember.isMaster()) {
								
								processApplicationParmRemove(message);
								
							}else if (message.getKey().equals(MessageKey.APPLICATION_PARMS_CLEAR_ALL) && !localTcpMember.isMaster()) {
								
								processApplicationParmsClearAll(message);
								
							}
							
							//session
							else if (message.getKey().equals(MessageKey.SESSION_NEW_USER) && !localTcpMember.isMaster()) {
								
								processSessionNewUser(message);
								
							}else if (message.getKey().equals(MessageKey.SESSION_REMOVE_USER) && !localTcpMember.isMaster()) {
								
								processSessionRemoveUser(message);
								
							}else if (message.getKey().equals(MessageKey.SESSION_PARM_PUT) && !localTcpMember.isMaster()) {
								
								processSessionParmPut(message);
								
							}else if (message.getKey().equals(MessageKey.SESSION_PARM_REMOVE) && !localTcpMember.isMaster()) {
								
								processSessionParmRemove(message);
								
							}else if (message.getKey().equals(MessageKey.SESSION_PARMS_CLEAR_ALL) && !localTcpMember.isMaster()) {
								
								processSessionParmsClearAll(message);
								
							}
							
							//exit
							else if(message.getKey().equals(MessageKey.EXIT)) {
								break;
							}
							
						}catch (Exception e) {
							logger.error("CoreMessage " + e.getMessage());
						}
						
						
					}

					mcSocket.leaveGroup(mcIPAddress);
					mcSocket.close();

				} catch (Exception e) {
					logger.error("multiCastListen " + e.getMessage());
				}
				
			}
		});
		
		thread.start();
	}
	
	private void processNewMemberJoin(CoreMessage message, ObjectMapper objectMapper, MulticastSocket mcSocket, DatagramPacket packet){
		
		try{
			InetAddress receivedAddressUDP = packet.getAddress();
			int receivedPortUDP = packet.getPort();
			Member receivedTCPMember = message.getMember();

			CoreMessage replyMessage = new CoreMessage();
			replyMessage.setKey(MessageKey.WELCOME_NEW_SERVER_MEMBER);
			replyMessage.setMembers(tcpMembers);

			byte[] replyMegBytes = objectMapper.writeValueAsString(
					replyMessage).getBytes();

			DatagramPacket replyPacket = new DatagramPacket(replyMegBytes, replyMegBytes.length);
			replyPacket.setAddress(receivedAddressUDP);
			replyPacket.setPort(receivedPortUDP);
			mcSocket.send(replyPacket);
			
			//add received address
			if(!receivedTCPMember.isMaster()){
				synchronized (tcpMembers) {
					tcpMembers.add(receivedTCPMember);
				}
				logger.info("new momber join => " + receivedTCPMember);
			}else{
				logger.info("master join => " + receivedTCPMember);
			}
			
			printTCPMembers();
		}catch(Exception e){
			logger.error("processNewMemberJoin " + e.getMessage());
		}
		
	}
	
	private void processRequestLiveData(CoreMessage message, ObjectMapper objectMapper, MulticastSocket mcSocket, DatagramPacket packet){
		
		try{
			InetAddress receivedAddressUDP = packet.getAddress();
			int receivedPortUDP = packet.getPort();

			CoreMessage replyMessage = new CoreMessage();
			replyMessage.setKey(MessageKey.LIVE_DATA);
			
			replyMessage.setSessionParm(sessionStore.getUsers());
			replyMessage.setApplicationContextParms(applicationContextData.getParms());

			byte[] replyMegBytes = objectMapper.writeValueAsString(
					replyMessage).getBytes();

			DatagramPacket replyPacket = new DatagramPacket(replyMegBytes, replyMegBytes.length);
			replyPacket.setAddress(receivedAddressUDP);
			replyPacket.setPort(receivedPortUDP);
			mcSocket.send(replyPacket);
			
		}catch(Exception e){
			logger.error("processRequestLiveData " + e.getMessage());
		}
		
	}
	
	private void processApplicationParmPut(CoreMessage message){
		
		if(applicationContextData == null){
			logger.error("applicationContextData null");
			return;
		}
		
		try{
			applicationContextData.getParms().putAll(message.getApplicationContextParms());
			
		}catch(Exception e){
			logger.error("processApplicationParmPut " + e.getMessage());
		}
		
	}
	
	private void processApplicationParmRemove(CoreMessage message){
		
		if(applicationContextData == null){
			logger.error("applicationContextData null");
			return;
		}
		
		try{
			for(String key : message.getApplicationContextParms().keySet()){
				applicationContextData.getParms().remove(key);
			}
		}catch(Exception e){
			logger.error("processApplicationParmRemove " + e.getMessage());
		}
		
	}
	
	private void processApplicationParmsClearAll(CoreMessage message){
		
		if(applicationContextData == null){
			logger.error("applicationContextData null");
			return;
		}
		
		try{
			applicationContextData.getParms().clear();
		}catch(Exception e){
			logger.error("processApplicationParmsClearAll " + e.getMessage());
		}
		
	}
	
	private void processSessionNewUser(CoreMessage message){
		
		if(sessionStore == null){
			logger.error("sessionStore null");
			return;
		}
		
		try{
			String sessionId = null;
			
			for(String key : message.getSessionParm().keySet()){
				sessionId = key;
			}
			
			if(!sessionStore.getUsers().containsKey(sessionId)){
				sessionStore.getUsers().put(sessionId, new HashMap<>());
			}
		}catch(Exception e){
			logger.error("processSessionNewUser " + e.getMessage());
		}
		
	}
	
	private void processSessionRemoveUser(CoreMessage message){
		
		if(sessionStore == null){
			logger.error("sessionStore null");
			return;
		}
		
		try{
			String sessionId = null;
			
			for(String key : message.getSessionParm().keySet()){
				sessionId = key;
			}
			
			if(sessionStore.getUsers().containsKey(sessionId)){
				sessionStore.getUsers().remove(sessionId);
			}
		}catch(Exception e){
			logger.error("processSessionRemoveUser " + e.getMessage());
		}
		
		
	}
	
	private void processSessionParmPut(CoreMessage message){
		
		if(sessionStore == null){
			logger.error("sessionStore null");
			return;
		}
		
		try{
			String sessionId = null;
			Map<String, Object> sessionPair = null;
			
			for(String key : message.getSessionParm().keySet()){
				sessionId = key;
				sessionPair = message.getSessionParm().get(key);
			}
			
			if(sessionStore.getUsers().containsKey(sessionId)){
				for(String parmKey : sessionPair.keySet()){
					sessionStore.getUsers().get(sessionId).put(parmKey, sessionPair.get(parmKey));
				}
			}
		}catch(Exception e){
			logger.error("processSessionAddParm " + e.getMessage());
		}
		
	}
	
	private void processSessionParmRemove(CoreMessage message){
		
		if(sessionStore == null){
			logger.error("sessionStore null");
			return;
		}
		
		try{
			String sessionId = null;
			Map<String, Object> sessionPair = null;
			
			for(String key : message.getSessionParm().keySet()){
				sessionId = key;
				sessionPair = message.getSessionParm().get(key);
				break;
			}
			
			if(sessionStore.getUsers().containsKey(sessionId)){
				for(String parmKey : sessionPair.keySet()){
					
					if(sessionStore.getUsers().get(sessionId).containsKey(parmKey)){
						sessionStore.getUsers().get(sessionId).remove(parmKey);
					}
					
				}
			}
		}catch(Exception e){
			logger.error("processSessionRemoveParm " + e.getMessage());
		}
		
	}
	
	private void processSessionParmsClearAll(CoreMessage message){
		
		if(sessionStore == null){
			logger.error("sessionStore null");
			return;
		}
		
		try{
			String sessionId = null;
			
			for(String key : message.getSessionParm().keySet()){
				sessionId = key;
			}
			
			if(sessionStore.getUsers().containsKey(sessionId))
				sessionStore.getUsers().get(sessionId).clear();
			
		}catch(Exception e){
			logger.error("processSessionClear " + e.getMessage());
		}
		
	}
	
	public List<Member> getTcpMembers() {
		return tcpMembers;
	}
	
	public void removeTcpMemberByMember(Member member){
		synchronized (tcpMembers) {
			tcpMembers.remove(member);
		}
		logger.info(Message.MEMBER_LEAVE.getMessage() + member);
		printTCPMembers();
	}
	
	public Member getLocalTcpMember(){
		return localTcpMember;
	}

	private void printTCPMembers(){
		
		if(tcpMembers.size() == 0){
			System.out.println("\n### EMPTY MEMBERS ###\n");
			return;
		}
		
		System.out.println("\n### MEMBERS ###");
		for(Member member : tcpMembers){
			
			StringBuilder print = new StringBuilder(member.toString());
			
			if(member.getPort() == localTcpMember.getPort())
				print.append(" [this]");
			if(member.isMaster())
				print.append(" [master]");
			
			System.out.println(print);
			
		}
		System.out.println("### MEMBERS ###\n");
	}
	
	public MulticastGroup getMulticastGroup(){
		return new MulticastGroup(multicastAddress, multicastPort);
	}

	public void setSessionStore(SessionStore sessionStore) {
		this.sessionStore = sessionStore;
	}
	
	public void setApplicationContextData(ApplicationContextData applicationContextData){
		this.applicationContextData = applicationContextData;
	}
}
