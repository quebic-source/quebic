package com.lovi.quebic.exception;

public class RequestProcessingException extends Exception{

	private static final long serialVersionUID = 400173508861017252L;
	
	private int statusCode = 400;

	public RequestProcessingException(String errorMessage) {
		super(errorMessage);
	}
	
	public int getStatusCode() {
		return statusCode;
	}

	public void setStatusCode(int statusCode) {
		this.statusCode = statusCode;
	}
	
	
}
