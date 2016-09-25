package com.lovi.quebic.sockm.config;

public class MulticastNetworkGroup {

	private String multicastAddress;
	private int multicastPort;
	
	public MulticastNetworkGroup(String multicastAddress, int multicastPort) {
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

	@Override
	public int hashCode() {
		final int prime = 31;
		int result = 1;
		result = prime * result + ((multicastAddress == null) ? 0 : multicastAddress.hashCode());
		result = prime * result + multicastPort;
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
		MulticastNetworkGroup other = (MulticastNetworkGroup) obj;
		if (multicastAddress == null) {
			if (other.multicastAddress != null)
				return false;
		} else if (!multicastAddress.equals(other.multicastAddress))
			return false;
		if (multicastPort != other.multicastPort)
			return false;
		return true;
	}
	
}
