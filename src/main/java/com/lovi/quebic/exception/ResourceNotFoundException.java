package com.lovi.quebic.exception;

public class ResourceNotFoundException extends Exception{

	private static final long serialVersionUID = -4099067435975569635L;

	public ResourceNotFoundException() {
		super(ErrorMessage.RESOURCE_NOT_FOUND.getMessage());
	}
	
	public ResourceNotFoundException(String message) {
		super(message);
	}
}
