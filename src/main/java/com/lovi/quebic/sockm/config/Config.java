package com.lovi.quebic.sockm.config;

import com.lovi.quebic.sockm.config.option.DiscoveryOption;
import com.lovi.quebic.sockm.config.option.TcpServerOption;

public class Config {
	
	private DiscoveryOption discoveryOption;
	private TcpServerOption tcpServerOption;

	public DiscoveryOption getDiscoveryOption() {
		return discoveryOption;
	}
	public void setDiscoveryOption(DiscoveryOption discoveryOption) {
		this.discoveryOption = discoveryOption;
	}
	public TcpServerOption getTcpServerOption() {
		return tcpServerOption;
	}
	public void setTcpServerOption(TcpServerOption tcpServerOption) {
		this.tcpServerOption = tcpServerOption;
	}
}
