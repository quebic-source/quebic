package com.lovi.quebic.web.impl;

import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.web.Request;
import com.lovi.quebic.web.enums.HttpMethod;
import com.lovi.quebic.web.file.HttpFile;

/**
 * 
 * @author Tharanga Thennakoon
 *
 */
public class RequestImpl implements Request {

	final static Logger logger = LoggerFactory.getLogger(Request.class);

	private String method;
	private String location;
	private String queryParameterString;
	private String protocolVersion;
	private Map<String, String> headers = new HashMap<>();
	private Map<String, String> parameters = new HashMap<>();
	private Map<String, HttpFile> files = new HashMap<>();

	public RequestImpl() throws Exception {
	}
	
	@Override
	public String getMethodStr() {
		return method;
	}

	@Override
	public String getLocation() {
		return location;
	}

	@Override
	public Map<String, String> getHeaders() {
		return headers;
	}

	@Override
	public String getHeader(String key) {
		return headers.get(key);
	}

	@Override
	public Map<String, String> getParameters() {
		return parameters;
	}

	@Override
	public String getParameter(String key) {
		return parameters.get(key);
	}

	@Override
	public String getContentType() {
		return getHeader("Content-Type");
	}
	
	@Override
	public String getContentLength() {
		return getHeader("Content-Length");
	}

	@Override
	public HttpFile getFile(String attributeName) {
		return files.get(attributeName);
	}
	
	@Override
	public Map<String, HttpFile> getFiles() {
		return files;
	}

	@Override
	public String getVersion() {
		return protocolVersion;
	}

	@Override
	public String getQueryParameterString() {
		return queryParameterString;
	}

	@Override
	public HttpMethod getMethod() {
		return HttpMethod.lookup(method);
	}

	@Override
	public InputStream getInputStream() {
		return null;
	}

	public void setMethodStr(String methodStr) {
		this.method = methodStr;
	}

	public void setLocation(String location) {
		this.location = location;
	}

	public void setHeaders(Map<String, String> headers) {
		this.headers = headers;
	}


	public void setParameters(Map<String, String> parameters) {
		this.parameters = parameters;
	}


	public void setVersion(String version) {
		this.protocolVersion = version;
	}

	public void setFiles(Map<String, HttpFile> files) {
		this.files = files;
	}
}