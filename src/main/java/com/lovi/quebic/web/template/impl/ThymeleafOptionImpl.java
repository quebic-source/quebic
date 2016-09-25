package com.lovi.quebic.web.template.impl;

import com.lovi.quebic.web.template.TemplateOption;

public class ThymeleafOptionImpl implements TemplateOption {

	private String mode;
	private String urlAccess;
	private String suffix;
	private String location;
	private boolean cacheable;
	
	/**
	 * 	<b>Default Template Options </b>
	 * 	<br/>
	  	urlAccess = /web/templates/[templateName]
	  	<br/>
	    mode = XHTML
	    <br/>
		suffix = .html
		<br/>
		location = web/templates
	 */
	public ThymeleafOptionImpl() {
		this.urlAccess = "/web/templates";
		this.mode = "LEGACYHTML5";
		this.suffix = ".html";
		this.location = "web/templates";
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
	public String getMode() {
		return mode;
	}

	@Override
	public void setMode(String mode) {
		this.mode = mode;
	}

	@Override
	public String getSuffix() {
		return suffix;
	}

	@Override
	public void setSuffix(String suffix) {
		this.suffix = suffix;
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
