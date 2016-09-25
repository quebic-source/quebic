package com.lovi.quebic.cluster.option;

import com.lovi.quebic.cluster.loadblancer.LoadBlancer;
import com.lovi.quebic.cluster.option.impl.ClusterOptionImpl;
import com.lovi.quebic.web.HttpServerOption;

public interface ClusterOption extends HttpServerOption{

	static ClusterOption createDefaultOption(){
		return new ClusterOptionImpl()
				.setLoadBlancer(LoadBlancer.createDefaultLoadBlancer())
				.setMulticastGroup(new MulticastGroup("230.1.1.1", 12345));
	}
	
	ClusterOption setMaster();

	boolean isMaster();
	
	ClusterOption setLoadBlancer(LoadBlancer loadBlancer);
	
	LoadBlancer getLoadBlancer();

	MulticastGroup getMulticastGroup();

	ClusterOption setMulticastGroup(MulticastGroup multicastGroup);
	
}
