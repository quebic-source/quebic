package com.lovi.quebic.cluster;

public class Member {

	private String address;
	private int port;
	private boolean master;
	
	public Member() {
	}
	
	public Member(String address, int port) {
		this.address = address;
		this.port = port;
	}
	
	public Member(String address, int port, boolean master) {
		this.address = address;
		this.port = port;
		this.master = master;
	}


	public String getAddress() {
		return address;
	}
	public void setAddress(String address) {
		this.address = address;
	}
	public int getPort() {
		return port;
	}
	public void setPort(int port) {
		this.port = port;
	}
	
	public boolean isMaster() {
		return master;
	}

	public void setMaster(boolean master) {
		this.master = master;
	}

	@Override
	public int hashCode() {
		final int prime = 31;
		int result = 1;
		result = prime * result + port;
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
		Member other = (Member) obj;
		if (port != other.port)
			return false;
		return true;
	}

	@Override
	public String toString() {
		return "Member [address=" + address + ", port=" + port + ", master="
				+ master + "]";
	}

}
