package com.lovi.quebic.web.file.impl;

import io.netty.handler.codec.http.multipart.FileUpload;

import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.async.AsyncExecutor;
import com.lovi.quebic.async.Future;
import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.web.file.HttpFile;

public class HttpFileAttributeImpl implements HttpFile {

	final Logger logger = LoggerFactory.getLogger(HttpFile.class);
	private StringBuilder resourcesSaveLocation = new StringBuilder("src/main/resources");
	private ByteBuffer byteBuffer;
	
	private String fieldName;
	private String fileName;
	private String fileContentType;
	
	private String fileExtention;

	public HttpFileAttributeImpl(FileUpload fileUpload) throws IOException {
		
		this.byteBuffer = fileUpload.getByteBuf().nioBuffer();
		
		this.fileName = fileUpload.getFilename();
		this.fieldName = fileUpload.getName();
		this.fileContentType = fileUpload.getContentType();
		
		this.fileExtention = fileName.substring(fileName.indexOf(".") + 1);
		
	}

	@Override
	public void save() {
		save("", null, null);
	}
	
	@Override
	public void save(String path) {
		save(path, null, null);
	}
	
	@Override
	public void save(Handler<String> successHandler, Handler<Throwable> failureHandler) {
		save("", successHandler, failureHandler);
	}

	@Override
	public void save(String path, Handler<String> successHandler, Handler<Throwable> failureHandler) {

		AsyncExecutor<Future<String>> asyncExecutor = AsyncExecutor.create();
		asyncExecutor.run(()->{
			
			Future<String> future = Future.create();
			
			try {
				
				String path_ = path;
				
				if(path_.startsWith("/"))
					path_ = path_.substring(1);
				
				if(path_.endsWith("/"))
					path_ = path_.substring(0, path_.length() - 1);
				
				
				if(path_.equals(""))
					resourcesSaveLocation.append("/").append(fileName);
				else
					resourcesSaveLocation.append("/").append(path_).append("/").append(fileName);
				
					
				String fileSavedLocation = resourcesSaveLocation.toString();
				
				logger.info("start-saving file > " + fileName);
				
				//write file
				Path file = Paths.get(fileSavedLocation);
				Files.write(file, byteBuffer.array());
				
				logger.info("end-saving file > " + fileName);
				
				future.setResult(fileSavedLocation);
				
				
			} catch (Exception e) {
				
				future.setFailure(e);
				
			}
			return future;
			
		}, future->{
			

			if(future.isSuccess()){
				
				if(successHandler != null)
					successHandler.handle(future.getResult());
			}else{
				
				if(failureHandler != null)
					failureHandler.handle(future.getFailure());
				
			}
			
			
		}, f->{
		
			if(failureHandler != null)
				failureHandler.handle(f);
			
		});
		
	}

	@Override
	public String getFieldName() {
		return fieldName;
	}

	@Override
	public String getContentType() {
		return fileContentType;
	}

	@Override
	public String getFileName() {
		return fileName;
	}

	@Override
	public HttpFile setFileName(String fileName) {
		this.fileName = fileName + "." + fileExtention;
		return this;
	}

	@Override
	public Integer getFileLength() {
		return byteBuffer.limit();
	}

	@Override
	public ByteBuffer getByteBuffer() {
		return byteBuffer;
	}

}
