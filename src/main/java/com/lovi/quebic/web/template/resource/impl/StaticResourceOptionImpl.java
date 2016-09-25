package com.lovi.quebic.web.template.resource.impl;

import com.lovi.quebic.web.template.resource.StaticResourceOption;

public class StaticResourceOptionImpl implements StaticResourceOption {

	private String urlAccess;
	private String location;
	private boolean cacheable;
	
	/**
	 * 	<b>Default Static Resources Options</b>
	 *  <br/>
	 *  urlAccess = /static/[resourceName]
	 *  <br/>
		location = web/static
	 */
	public StaticResourceOptionImpl() {
		this.urlAccess = "/static";
		this.location = "web/static";
	}
	
	@Override
	public String getUrlAccess() {
		return urlAccess;
	}

	@Override
	public void setUrlAccess(String urlAccess) {
		this.urlAccess = urlAccess;
	}

	@Override
	public String getLocation() {
		return location + "/";
	}

	@Override
	public void setLocation(String location) {
		this.location = location;
	}

	@Override
	public void setCacheable(boolean cacheable) {
		this.cacheable = cacheable;
	}
	
	@Override
	public boolean getCacheable() {
		return this.cacheable;
	}

}
