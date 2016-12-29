package com.lovi.quebic.web.impl;

import java.util.HashSet;
import java.util.Set;

import com.lovi.quebic.handlers.FailureHandler;
import com.lovi.quebic.web.RequestMap;
import com.lovi.quebic.web.RequestMapper;
import com.lovi.quebic.web.ServerContext;
import com.lovi.quebic.web.enums.HttpMethod;

public class RequestMapperImpl implements RequestMapper{

	private Set<RequestMap> requestMaps = new HashSet<>();
	private FailureHandler<ServerContext> failureHandler;

	@Override
	public RequestMap map(String path, HttpMethod method){
		return new RequestMapImpl(this, path, method);
	}

	@Override
	public Set<RequestMap> getRequestMaps() {
		return requestMaps;
	}
	
	@Override
	public FailureHandler<ServerContext> getFailureHandler() {
		return failureHandler;
	}

	@Override
	public RequestMapperImpl setFailureHandler(FailureHandler<ServerContext> failureHandler) {
		this.failureHandler = failureHandler;
		return this;
	}
	
}
