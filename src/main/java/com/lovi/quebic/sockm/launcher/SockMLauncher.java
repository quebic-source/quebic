package com.lovi.quebic.sockm.launcher;

import java.util.List;
import java.util.Map;

import com.lovi.quebic.sockm.config.Config;
import com.lovi.quebic.sockm.config.MulticastNetworkGroup;
import com.lovi.quebic.sockm.config.TcpAddress;
import com.lovi.quebic.sockm.config.option.DiscoveryOption;
import com.lovi.quebic.sockm.launcher.impl.SockMLauncherImpl;
import com.lovi.quebic.sockm.messenger.Messenger;

public interface SockMLauncher {

	static SockMLauncher create(){
		return SockMLauncherImpl.create();
	}
	
	static SockMLauncher create(Config config){
		return SockMLauncherImpl.create(config);
	}
	
	<K, V> Map<K, V> getMap(String key);
	
	<K, V> Map<K, V> getMap(String key, Class<K> typeOfKey, Class<V> typeOfValue);

	<E> List<E> getList(String key);
	
	<E> List<E> getList(String key, Class<E> typeOfValue);
	
	Messenger getMessenger();
	
	/*
	ContextMapData getContextMapData();
	
	ContextListData getContextListData();
	*/
	void stop();
	
	String getLauncherId();
	
	MulticastNetworkGroup getMulticastGroup();
	
	/**
	 * Tcp Server Address from Config
	 * @return
	 */
	TcpAddress getTcpServerAddress();
	
	/**
	 * Running Tcp Server Address
	 * @return
	 */
	TcpAddress getRunningTcpServerAddress();

	DiscoveryOption getDiscoveryOption();

}
