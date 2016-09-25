package com.lovi.quebic.web;

import java.io.InputStream;
import java.util.Map;

import com.lovi.quebic.web.enums.HttpMethod;
import com.lovi.quebic.web.file.HttpFile;

/**
 * 
 * @author Tharanga Thennakoon
 *
 */
public interface Request {

	HttpMethod getMethod();
	
	String getMethodStr();

	String getLocation();

	Map<String,String> getHeaders();
	
	String getHeader(String key);
	
	Map<String, String> getParameters();
	
	String getParameter(String key);
	
	String getContentType();

	String getContentLength();
	
	HttpFile getFile(String fieldName);

	String getVersion();

	String getQueryParameterString();

	InputStream getInputStream();

	Map<String, HttpFile> getFiles();
	
}
