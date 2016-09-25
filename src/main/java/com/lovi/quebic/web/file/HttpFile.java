package com.lovi.quebic.web.file;

import java.nio.ByteBuffer;

import com.lovi.quebic.handlers.Handler;

public interface HttpFile{
	
	void save();
	void save(String path);
	void save(Handler<String> successHandler, Handler<Throwable> failureHandler);
	void save(String path, Handler<String> successHandler, Handler<Throwable> failureHandler);
	Integer getFileLength();
	ByteBuffer getByteBuffer();
	String getFileName();
	HttpFile setFileName(String fileName);
	String getFieldName();
	String getContentType();
	
}
