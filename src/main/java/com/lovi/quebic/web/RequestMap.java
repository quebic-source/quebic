package com.lovi.quebic.web;

import com.lovi.quebic.handlers.FailureHandler;
import com.lovi.quebic.handlers.RequestHandler;
import com.lovi.quebic.web.enums.HttpMethod;

public interface RequestMap{

	RequestHandler<ServerContext> getHandler();

	RequestMap setHandler(RequestHandler<ServerContext> handler);

	String getPath();
	
	String getRegExpPath();
	
	void setRegExpPath(String regExpPath);

	HttpMethod getHttpMethod();

	FailureHandler<ServerContext> getFailureHandler();

	RequestMap setFailureHandler(FailureHandler<ServerContext> failureHandler);

}
