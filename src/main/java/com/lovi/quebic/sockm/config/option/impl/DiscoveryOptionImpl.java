package com.lovi.quebic.sockm.config.option.impl;

import java.util.HashSet;
import java.util.Set;

import com.lovi.quebic.sockm.config.MulticastNetworkGroup;
import com.lovi.quebic.sockm.config.option.DiscoveryOption;

public class DiscoveryOptionImpl implements DiscoveryOption {

	private boolean enableMulticastMethod;
	private boolean enableTcpMethod;
	
	private Set<MulticastNetworkGroup> tcpMembers = new HashSet<>();
	private MulticastNetworkGroup multicastGroup;
	
	@Override
	public DiscoveryOption enableMulticastMethod(){
		enableMulticastMethod = true;
		enableTcpMethod = false;
		return this;
	}
	
	@Override
	public DiscoveryOption enableTcpMethod(){
		enableTcpMethod = true;
		enableMulticastMethod = false;
		return this;
	}
	
	@Override
	public boolean isEnableMulticastMethod(){
		return enableMulticastMethod;
	}
	
	@Override
	public boolean isEnableTcpMethod(){
		return enableTcpMethod;
	}
	
	@Override
	public MulticastNetworkGroup getMulticastGroup() {
		return multicastGroup;
	}

	@Override
	public void setMulticastGroup(MulticastNetworkGroup multicastGroup) {
		this.multicastGroup = multicastGroup;
	}
	
	@Override
	public DiscoveryOption addTcpMember(MulticastNetworkGroup member){
		tcpMembers.add(member);
		return this;
	}
	
	@Override
	public Set<MulticastNetworkGroup> getTcpMembers() {
		return tcpMembers;
	}

	
	
}
