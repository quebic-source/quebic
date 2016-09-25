package com.lovi.quebic.exception;

public class InternalServerException extends Exception{
	
	private static final long serialVersionUID = -5801139621474961128L;
	private int statusCode = 500;

	public InternalServerException(String errorMessage) {
		super(errorMessage);
	}
	
	public int getStatusCode() {
		return statusCode;
	}

	public void setStatusCode(int statusCode) {
		this.statusCode = statusCode;
	}
}
