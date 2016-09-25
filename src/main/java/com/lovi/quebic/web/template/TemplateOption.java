package com.lovi.quebic.web.template;

import com.lovi.quebic.web.HttpServerOption;
import com.lovi.quebic.web.template.impl.ThymeleafOptionImpl;

public interface TemplateOption extends HttpServerOption{
	
	/**
	 * 	<b>Thymeleaf Template Options </b>
	 * 	<br/>
	  	urlAccess = /web/templates/[templateName]
	  	<br/>
	    mode = XHTML
	    <br/>
		suffix = .html
		<br/>
		location = web/templates
	 */
	static TemplateOption createThymeleafOption(){
		return new ThymeleafOptionImpl();
	}

	String getUrlAccess();

	void setUrlAccess(String urlAccess);
	
	String getMode();

	void setMode(String mode);

	String getSuffix();

	void setSuffix(String suffix);

	String getLocation();

	void setLocation(String location);
	
	void setCacheable(boolean cacheable);

	boolean getCacheable();
}
