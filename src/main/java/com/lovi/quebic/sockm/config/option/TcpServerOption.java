package com.lovi.quebic.sockm.config.option;

import com.lovi.quebic.sockm.config.TcpAddress;
import com.lovi.quebic.sockm.config.option.impl.TcpServerOptionImpl;


public interface TcpServerOption {

	static TcpServerOption create(){
		return new TcpServerOptionImpl();
	}

	TcpAddress getAddress();

	void setAddress(TcpAddress address);
	
}
