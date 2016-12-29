package com.lovi.quebic.web.impl;

import java.util.Map;
import java.util.Set;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.thymeleaf.context.Context;
import org.thymeleaf.context.VariablesMap;

import com.lovi.quebic.exception.ErrorMessage;
import com.lovi.quebic.web.ApplicationContextData;
import com.lovi.quebic.web.HttpServer;
import com.lovi.quebic.web.Request;
import com.lovi.quebic.web.Response;
import com.lovi.quebic.web.ServerContext;
import com.lovi.quebic.web.Session;
import com.lovi.quebic.web.enums.HttpMethod;
import com.lovi.quebic.web.server.HttpServerHandler;

public class ServerContextImpl implements ServerContext {

	final static Logger logger = LoggerFactory.getLogger(ServerContext.class);
	
	private HttpServer httpServer;
	private HttpServerHandler nettyHttpServerHandler;
	private Request httpRequst;
	private Response httpResponse;
	private ApplicationContextData applicationContextData;
	private Session session;
	private Context context;

	public ServerContextImpl(HttpServer httpServer, HttpServerHandler nettyHttpServerHandler, Request httpRequst,
			Response httpResponse, ApplicationContextData applicationContextData,Session session) {
		this.httpServer = httpServer;
		this.nettyHttpServerHandler = nettyHttpServerHandler;
		this.httpRequst = httpRequst;
		this.httpResponse = httpResponse;
		this.applicationContextData = applicationContextData;
		this.session = session;
		this.context = new Context();
	}

	@Override
	public void forward(String path) {
		forward(path, HttpMethod.GET);
	}
	
	@Override
	public void forward(String path, HttpMethod method) {
		try {
			nettyHttpServerHandler.requestProcess(this, path, method);
		} catch (Exception e) {
			logger.error(ErrorMessage.REQUEST_FORWORD_ERROR.getMessage());
		}

	}
	
	@Override
	public void redirect(String path){
		httpResponse.setResponseCode(301);
		httpResponse.setHeader("Location", path);
		httpResponse.write("");
	}
	
	@Override
	public Request getHttpRequst() {
		return httpRequst;
	}

	@Override
	public Response getHttpResponse() {
		return httpResponse;
	}

	@Override
	public ApplicationContextData getApplicationContextData() {
		return applicationContextData;
	}
	
	@Override
	public Session getSession() {
		return session;
	}
	
	@Override
	public Context getTemplateContext(){
		return context;
	}
	
	@Override
	public void putData(String key, Object value){
		context.setVariable(key, value);
	}
	
	@Override
	public Object getData(String key){
		//return context.getVariable(key);
		VariablesMap<String, Object> variablesMap = context.getVariables();
		return variablesMap.get(key);
	}
	
	@Override
	public void removeData(String key){
		//context.removeVariable(key);
		VariablesMap<String, Object> variablesMap = context.getVariables();
		variablesMap.remove(key);
	}
	
	@Override
	public void clearAllData(){
		context.clearVariables();
	}
	
	@Override
	public Set<String> getDataKeySet(){
		//return context.getVariableNames();
		VariablesMap<String, Object> variablesMap = context.getVariables();
		return variablesMap.keySet();
	}
	
	@Override
	public void setDataMap(Map<String, Object> data){
		context.setVariables(data);
	}
	
	@Override
	public boolean isContainsData(String key){
		//return context.containsVariable(key);
		VariablesMap<String, Object> variablesMap = context.getVariables();
		return variablesMap.containsKey(key);
	}

	@Override
	public void loadTemplate(String templateName){
		putData("templateName", templateName);
		forward(httpServer.getTemplateOption().getUrlAccess());
	}
}
