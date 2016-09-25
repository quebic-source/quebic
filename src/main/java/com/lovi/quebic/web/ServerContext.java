package com.lovi.quebic.web;

import java.util.Map;
import java.util.Set;

import org.thymeleaf.context.Context;

import com.lovi.quebic.web.enums.HttpMethod;

public interface ServerContext {
	
	Request getHttpRequst();

	Response getHttpResponse();

	ApplicationContextData getApplicationContextData();
	
	Session getSession();	

	/**
	 * this method is used for request forwarding. HTTP Method -> GET
	 * @param path
	 */
	void forward(String path);
	
	/**
	 * this method is used for request forwarding.
	 * @param path
	 * @param method
	 */
	void forward(String path, HttpMethod method);

	/**
	 * this method redirect the request to new location. only support for GET request
	 * @param path
	 */
	void redirect(String path);

	Context getTemplateContext();
	
	void putData(String key, Object value);

	Object getData(String key);

	void removeData(String key);

	void clearAllData();

	Set<String> getDataKeySet();

	void setDataMap(Map<String, Object> data);

	boolean isContainsData(String key);

	void loadTemplate(String templateName);

}
