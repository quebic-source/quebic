package com.lovi.quebic.sockm.messenger;

import com.lovi.quebic.sockm.config.TcpAddress;

public class SubscriberDirectory {
	
	/**
	 * tcpAddress is used only for Tcp Discovery mode
	 */
	private TcpAddress tcpAddress;
	private String launcherId;
	private String listenerId;
	private String listenerAddress;
	
	public SubscriberDirectory() {
	}
	
	public SubscriberDirectory(TcpAddress tcpAddress, String launcherId, String listenerId, String listenerAddress) {
		this.tcpAddress = tcpAddress;
		this.launcherId = launcherId;
		this.listenerId = listenerId;
		this.listenerAddress = listenerAddress;
	}
	
	public TcpAddress getTcpAddress() {
		return tcpAddress;
	}
	public void setTcpAddress(TcpAddress tcpAddress) {
		this.tcpAddress = tcpAddress;
	}
	public String getLauncherId() {
		return launcherId;
	}
	public void setLauncherId(String launcherId) {
		this.launcherId = launcherId;
	}
	public String getListenerId() {
		return listenerId;
	}
	public void setListenerId(String listenerId) {
		this.listenerId = listenerId;
	}
	public String getListenerAddress() {
		return listenerAddress;
	}
	public void setListenerAddress(String listenerAddress) {
		this.listenerAddress = listenerAddress;
	}
	@Override
	public String toString() {
		return "MessageListener [tcpAddress=" + tcpAddress + ", listenerId="
				+ listenerId + ", listenerAddress=" + listenerAddress + "]";
	}
	@Override
	public int hashCode() {
		final int prime = 31;
		int result = 1;
		result = prime * result
				+ ((listenerId == null) ? 0 : listenerId.hashCode());
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
		SubscriberDirectory other = (SubscriberDirectory) obj;
		if (listenerId == null) {
			if (other.listenerId != null)
				return false;
		} else if (!listenerId.equals(other.listenerId))
			return false;
		return true;
	}

}
