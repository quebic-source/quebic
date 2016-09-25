package com.lovi.quebic.exception;

public class RequestForwardException extends Exception{

	private static final long serialVersionUID = -4713030709320239765L;
	
	public RequestForwardException() {
		super(ErrorMessage.REQUEST_FORWORD_ERROR.getMessage());
	}
	
	public RequestForwardException(String message) {
		super(message);
	}

}
