package com.lovi.quebic.cluster.option;

public class MulticastGroup {

	private String multicastAddress;
	private int multicastPort;
	
	public MulticastGroup(String multicastAddress, int multicastPort) {
		this.multicastAddress = multicastAddress;
		this.multicastPort = multicastPort;
	}

	public String getMulticastAddress() {
		return multicastAddress;
	}

	public void setMulticastAddress(String multicastAddress) {
		this.multicastAddress = multicastAddress;
	}

	public int getMulticastPort() {
		return multicastPort;
	}

	public void setMulticastPort(int multicastPort) {
		this.multicastPort = multicastPort;
	}
	
	
}
