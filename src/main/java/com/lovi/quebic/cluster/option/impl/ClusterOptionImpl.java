package com.lovi.quebic.cluster.option.impl;

import com.lovi.quebic.cluster.loadblancer.LoadBlancer;
import com.lovi.quebic.cluster.option.ClusterOption;
import com.lovi.quebic.cluster.option.MulticastGroup;

public class ClusterOptionImpl implements ClusterOption {

	private boolean master;
	private LoadBlancer loadBlancer;
	private MulticastGroup multicastGroup;
	
	@Override
	public ClusterOption setMaster() {
		master = true;
		return this;
	}
	
	@Override
	public boolean isMaster(){
		return master;
	}

	@Override
	public LoadBlancer getLoadBlancer() {
		return loadBlancer;
	}

	@Override
	public ClusterOption setLoadBlancer(LoadBlancer loadBlancer) {
		this.loadBlancer = loadBlancer;
		return this;
	}

	@Override
	public MulticastGroup getMulticastGroup() {
		return multicastGroup;
	}

	@Override
	public ClusterOption setMulticastGroup(MulticastGroup multicastGroup) {
		this.multicastGroup = multicastGroup;
		return this;
	}

	
}
