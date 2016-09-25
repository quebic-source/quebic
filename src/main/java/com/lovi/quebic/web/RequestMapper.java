package com.lovi.quebic.web;

import java.util.Set;

import com.lovi.quebic.handlers.FailureHandler;
import com.lovi.quebic.web.enums.HttpMethod;
import com.lovi.quebic.web.impl.RequestMapperImpl;

public interface RequestMapper {

	static RequestMapper create(){
		return new RequestMapperImpl();
	}
	
	Set<RequestMap> getRequestMaps();

	RequestMap map(String path, HttpMethod method);

	FailureHandler<ServerContext> getFailureHandler();

	RequestMapperImpl setFailureHandler(FailureHandler<ServerContext> failureHandler);

}
