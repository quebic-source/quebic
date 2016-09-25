package com.lovi.quebic.sockm.config.option;

import java.util.Set;

import com.lovi.quebic.sockm.config.MulticastNetworkGroup;
import com.lovi.quebic.sockm.config.option.impl.DiscoveryOptionImpl;

public interface DiscoveryOption {

	static DiscoveryOption create(){
		return new DiscoveryOptionImpl();
	}
	
	DiscoveryOption enableMulticastMethod();
	
	DiscoveryOption enableTcpMethod();
	
	boolean isEnableMulticastMethod();

	boolean isEnableTcpMethod();
	
	MulticastNetworkGroup getMulticastGroup();

	void setMulticastGroup(MulticastNetworkGroup multicastGroup);

	DiscoveryOption addTcpMember(MulticastNetworkGroup member);

	Set<MulticastNetworkGroup> getTcpMembers();

}
