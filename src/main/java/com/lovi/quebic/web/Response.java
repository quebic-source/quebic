package com.lovi.quebic.web;

import java.util.Map;

import com.lovi.quebic.exception.JsonParserException;

public interface Response {

	int getResponseCode();

	String getResponseReason();

	void setResponseCode(int responseCode);

	void setResponseReason(String responseReason);

	String getHeader(String header);

	void setHeader(String key, String value);
	
	void setContentType(String contentType);

	Map<String, String> getHeaders();

	byte[] getContent();
	
	void setEncodeAsGzip(boolean encodeAsGzip);

	void setKeepAlive(boolean keepAlive);

	void write(String content);

	void write(Integer content);

	void write(Short content);

	void write(Long content);

	void write(Float content);

	void write(Double content);

	void write(Boolean content);

	void writePOJO(Object content) throws JsonParserException;
	
	void write(byte[] content);

	String getContentType();

	void setChunkedTransfer(boolean chunkedTransfer);

	boolean isChunkedTransfer();

	boolean isEncodeAsGzip();

	boolean isKeepAlive();


}
