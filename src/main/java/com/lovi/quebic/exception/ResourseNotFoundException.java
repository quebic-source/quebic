package com.lovi.quebic.exception;

public class ResourseNotFoundException extends Exception{
	
	private static final long serialVersionUID = -918451678352909786L;
	private int statusCode = 404;

	public ResourseNotFoundException(String errorMessage) {
		super(errorMessage);
	}
	
	public int getStatusCode() {
		return statusCode;
	}

	public void setStatusCode(int statusCode) {
		this.statusCode = statusCode;
	}
}
