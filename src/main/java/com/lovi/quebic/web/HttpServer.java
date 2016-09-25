package com.lovi.quebic.web;

import java.util.List;

import com.lovi.quebic.cluster.ClusterConnector;
import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.web.enums.HttpMethod;
import com.lovi.quebic.web.impl.HttpServerImpl;
import com.lovi.quebic.web.template.TemplateOption;

/**
 * 
 * @author Tharanga Thennakoon
 *
 */
public interface HttpServer {
	
	static HttpServer create(){
		return new HttpServerImpl();
	}
	
	void setRequestMapper(RequestMapper requestMapper);

	void run(int port, RequestMapper requestMapper);

	void run(String hostname, int port, RequestMapper requestMapper);

	void run(int port, RequestMapper requestMapper, Handler<String> successHandler, Handler<Throwable> failureHandler);

	void run(String hostname, int port, RequestMapper requestMapper, Handler<String> successHandler,
			Handler<Throwable> failureHandler);
	
	/**
	 * this method is used for request forwarding. HTTP Method -> GET
	 * @param serverContext
	 * @param newLocation
	 * @throws Exception
	 */
	void requestProcess(ServerContext serverContext, String newLocation)throws Exception;

	/**
	 * this method is used for request forwarding
	 * @param serverContext
	 * @throws Exception
	 */
	void requestProcess(ServerContext serverContext, String newLocation, HttpMethod httpMethod) throws Exception;

	/**
	 * add HttpServerOption
	 * @param option
	 */
	void addHttpServerOption(HttpServerOption option);

	List<HttpServerOption> getHttpServerOption();
	
	TemplateOption getTemplateOption();

	RequestMapper getRequestMapper();
	
	ClusterConnector getClusterConnector();
	
	public SessionStore getSessionStore();

	public ApplicationContextData getApplicationContextData();
	

}
