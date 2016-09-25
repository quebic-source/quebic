package com.lovi.quebic.sockm.launcher.impl;

import java.util.List;
import java.util.Map;
import java.util.UUID;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.sockm.config.Config;
import com.lovi.quebic.sockm.config.MulticastNetworkGroup;
import com.lovi.quebic.sockm.config.TcpAddress;
import com.lovi.quebic.sockm.config.option.DiscoveryOption;
import com.lovi.quebic.sockm.config.option.TcpServerOption;
import com.lovi.quebic.sockm.context.ContextListData;
import com.lovi.quebic.sockm.context.ContextMapData;
import com.lovi.quebic.sockm.context.impl.ContextListDataImpl;
import com.lovi.quebic.sockm.context.impl.ContextMapDataImpl;
import com.lovi.quebic.sockm.exception.MulticastServerException;
import com.lovi.quebic.sockm.exception.TcpServerException;
import com.lovi.quebic.sockm.launcher.SockMLauncher;
import com.lovi.quebic.sockm.log.error.ErrorMessage;
import com.lovi.quebic.sockm.messenger.Messenger;
import com.lovi.quebic.sockm.server.MulticastServer;
import com.lovi.quebic.sockm.server.TcpServer;

public class SockMLauncherImpl implements SockMLauncher {

	private final Logger logger = LoggerFactory.getLogger(SockMLauncher.class);

	private final String id;

	private static SockMLauncher instance;
	private static Object lock = new Object();

	private Config config;

	private ContextMapData contextMapData;
	private ContextListData contextListData;
	
	private Messenger messenger;

	private MulticastServer multicastServer;
	private TcpServer tcpServer;

	private SockMLauncherImpl() {
		this(new Config());
	}

	private SockMLauncherImpl(Config config) {
		
		id = genarateId();
		
		this.config = config;

		prepareDiscoveryOption(config.getDiscoveryOption());
		prepareTcpServerOption(config.getTcpServerOption());

		prepareContextData();
		prepareMessenger();
		
		prepareMulticastServer();

		//prepareTcpServer();

	}

	private void prepareDiscoveryOption(DiscoveryOption option) {
		if (option == null) {
			// default options
			DiscoveryOption discoveryOption = DiscoveryOption.create();
			discoveryOption.enableMulticastMethod();
			discoveryOption.setMulticastGroup(new MulticastNetworkGroup("230.1.1.1",1234));
			
			config.setDiscoveryOption(discoveryOption);
		} 

	}
	
	private void prepareTcpServerOption(TcpServerOption option) {
		if(option == null){
			//default options
			TcpServerOption tcpServerOption = TcpServerOption.create();
			config.setTcpServerOption(tcpServerOption);
		}
		
	}

	private void prepareContextData() {
		contextMapData = new ContextMapDataImpl(this);
		contextListData = new ContextListDataImpl(this);
	}
	
	private void prepareMessenger() {
		messenger = Messenger.create(this);
	}

	private void prepareMulticastServer() {
		multicastServer = new MulticastServer(this, contextMapData, contextListData, messenger);
		try {
			multicastServer.listen();
		} catch (MulticastServerException e) {
			logger.error(ErrorMessage.MULTICAST_SERVER_ERROR + " {}",
					e.getMessage());
		}
	}
	
	private void prepareTcpServer(){
		tcpServer = new TcpServer(this);
		try {
			tcpServer.listen();
		} catch (TcpServerException e) {
			logger.error(ErrorMessage.TCP_SERVER_ERROR + " {}",
					e.getMessage());
		}
	}

	public static SockMLauncher create() {

		synchronized (lock) {
			if(instance == null)
				instance = new SockMLauncherImpl();
		}

		return instance;
	}

	public static SockMLauncher create(Config config) {

		synchronized (lock) {
			if(instance == null)
				instance = new SockMLauncherImpl();
		}

		return instance;
	}

	@Override
	public MulticastNetworkGroup getMulticastGroup() {
		return config.getDiscoveryOption().getMulticastGroup();
	}

	@Override
	public TcpAddress getTcpServerAddress(){
		return config.getTcpServerOption().getAddress();
		
	}
	
	@Override
	public TcpAddress getRunningTcpServerAddress(){
		return tcpServer.getLocalTcpAddress();
		
	}
	
	@Override
	public String getLauncherId() {
		return this.id;
	}

	@Override
	public <K, V> Map<K, V> getMap(String key) {
		return contextMapData.getMap(key);
	}

	@Override
	public <K, V> Map<K, V> getMap(String key, Class<K> typeOfKey,
			Class<V> typeOfValue) {
		return contextMapData.getMap(key, typeOfKey, typeOfValue);
	}

	@Override
	public <E> List<E> getList(String key) {
		return contextListData.getList(key);
	}

	@Override
	public <E> List<E> getList(String key, Class<E> typeOfValue) {
		return contextListData.getList(key, typeOfValue);
	}

	@Override
	public Messenger getMessenger() {
		return messenger;
	}

	@Override
	public void stop() {
		if(multicastServer != null)
			multicastServer.stopServer();
		
		if(tcpServer != null)
			tcpServer.stopServer();
	}

	@Override
	public DiscoveryOption getDiscoveryOption() {
		return config.getDiscoveryOption();
	}

	private String genarateId() {
		return UUID.randomUUID().toString();
	}

	@Override
	public int hashCode() {
		final int prime = 31;
		int result = 1;
		result = prime * result + ((id == null) ? 0 : id.hashCode());
		return result;
	}

	@Override
	public boolean equals(Object obj) {
		if (this == obj)
			return true;
		if (obj == null)
			return false;
		if (getClass() != obj.getClass())
			return false;
		SockMLauncherImpl other = (SockMLauncherImpl) obj;
		if (id == null) {
			if (other.id != null)
				return false;
		} else if (!id.equals(other.id))
			return false;
		return true;
	}
	
}
