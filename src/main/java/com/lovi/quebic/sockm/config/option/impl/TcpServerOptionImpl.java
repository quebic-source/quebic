package com.lovi.quebic.sockm.config.option.impl;

import com.lovi.quebic.sockm.config.TcpAddress;
import com.lovi.quebic.sockm.config.option.TcpServerOption;

public class TcpServerOptionImpl implements TcpServerOption {

	private TcpAddress address;

	public TcpServerOptionImpl() {
	}

	@Override
	public TcpAddress getAddress() {
		return address;
	}

	@Override
	public void setAddress(TcpAddress address) {
		this.address = address;
	}
	
	
}
