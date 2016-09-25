package com.lovi.quebic.context;

import org.springframework.context.annotation.Configuration;

import com.lovi.quebic.sockm.launcher.SockMLauncher;
import com.lovi.quebic.web.HttpServer;

@Configuration
public class AppConfig {
	
	private Class<?> appClass;
	private HttpServer httpServer;
	private SockMLauncher sockMLauncher;
	
	public Class<?> getAppClass() {
		return appClass;
	}
	public void setAppClass(Class<?> appClass) {
		this.appClass = appClass;
	}
	public HttpServer getHttpServer() {
		return httpServer;
	}
	public void setHttpServer(HttpServer httpServer) {
		this.httpServer = httpServer;
	}
	public SockMLauncher getSockMLauncher() {
		return sockMLauncher;
	}
	public void setSockMLauncher(SockMLauncher sockMLauncher) {
		this.sockMLauncher = sockMLauncher;
	}
}
