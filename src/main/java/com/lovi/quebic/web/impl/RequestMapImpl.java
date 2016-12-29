package com.lovi.quebic.web.impl;

import com.lovi.quebic.handlers.FailureHandler;
import com.lovi.quebic.handlers.RequestHandler;
import com.lovi.quebic.web.RequestMap;
import com.lovi.quebic.web.RequestMapper;
import com.lovi.quebic.web.ServerContext;
import com.lovi.quebic.web.enums.HttpMethod;

public class RequestMapImpl implements RequestMap{

	private String path;
	private String regExpPath;
	private HttpMethod httpMethod;
	private RequestHandler<ServerContext> handler;
	private FailureHandler<ServerContext> failureHandler;
	
	public RequestMapImpl(RequestMapper requestMapper, String path, HttpMethod httpMethod) {
		requestMapper.getRequestMaps().add(this);
		this.path = path;
		this.httpMethod = httpMethod;
	}

	@Override
	public RequestHandler<ServerContext> getHandler() {
		return handler;
	}

	@Override
	public RequestMap setHandler(RequestHandler<ServerContext> handler) {
		this.handler = handler;
		return this;
	}
	
	@Override
	public String getPath() {
		return path;
	}

	@Override
	public HttpMethod getHttpMethod() {
		return httpMethod;
	}

	@Override
	public String getRegExpPath() {
		return regExpPath;
	}

	@Override
	public void setRegExpPath(String regExpPath) {
		this.regExpPath = regExpPath;
	}

	@Override
	public FailureHandler<ServerContext> getFailureHandler() {
		return failureHandler;
	}

	@Override
	public RequestMap setFailureHandler(FailureHandler<ServerContext> failureHandler) {
		this.failureHandler = failureHandler;
		return this;
	}
	
}
