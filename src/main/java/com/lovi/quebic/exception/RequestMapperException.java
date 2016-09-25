package com.lovi.quebic.exception;

public class RequestMapperException extends Exception{

	private static final long serialVersionUID = -5513738908441171023L;

	public RequestMapperException() {
		super(ErrorMessage.UNABLE_TO_FOUND_REQUEST_MAPPER.getMessage());
	}
	
	public RequestMapperException(String message) {
		super(message);
	}
}
