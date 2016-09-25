package com.lovi.quebic.web.option;

import com.lovi.quebic.web.HttpServerOption;
import com.lovi.quebic.web.option.impl.SSLServerConfigImpl;

public interface SSLServerConfig extends HttpServerOption{

	static SSLServerConfig createDefaultOption(){
		return new SSLServerConfigImpl();
	}
	
}
