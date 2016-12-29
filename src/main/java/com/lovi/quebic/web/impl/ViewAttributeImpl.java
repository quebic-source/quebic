package com.lovi.quebic.web.impl;

import com.lovi.quebic.web.ServerContext;
import com.lovi.quebic.web.ViewAttribute;

public class ViewAttributeImpl implements ViewAttribute {

	private ServerContext serverContext;
	
	public ViewAttributeImpl(ServerContext serverContext) {
		this.serverContext = serverContext;
	}
	
	@Override
	public void put(String key, Object object) {
		serverContext.putData(key, object);
	}

	@Override
	public Object get(String key) {
		return serverContext.getData(key);
	}

	@Override
	public void loadView(String name) {
		serverContext.loadTemplate(name);
	}

}
