package com.lovi.quebic.web.template.resource;

import com.lovi.quebic.web.HttpServerOption;
import com.lovi.quebic.web.template.resource.impl.StaticResourceOptionImpl;

public interface StaticResourceOption extends HttpServerOption{

	static StaticResourceOption create(){
		return new StaticResourceOptionImpl();
	}
	
	String getUrlAccess();

	void setUrlAccess(String urlAccess);
	
	String getLocation();

	void setLocation(String location);
	
	void setCacheable(boolean cacheable);

	boolean getCacheable();
}
