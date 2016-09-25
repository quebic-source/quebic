package com.lovi.quebic.async.impl;

import com.lovi.quebic.async.Future;
import com.lovi.quebic.handlers.Handler;

/**
 * 
 * @author Tharanga Thennakoon
 *
 * @param <T>
 */
public class FutureImpl<T> implements Future<T> {
	
	private Handler<T> successHandler;
	private Handler<Throwable> failureHandler;
	private T result;
	private Throwable failure;
	
	private boolean isCompleted;
	private boolean success = true;
	
	private Integer errorCode = null;

	@Override
	public void setResult(T result) {
		this.result = result;
		this.success = true;
		this.isCompleted = true;
		if(successHandler != null){
			try {
				successHandler.handle(result);
			} catch (Exception e) {
				setFailure(e);
			}
		}
	}
	
	@Override
	public T getResult() {
		return result;
	}

	@Override
	public void setFailure(Throwable failure) {
		this.failure = failure;
		this.success = false;
		this.isCompleted = true;
		if(failureHandler != null){
			try {
				failureHandler.handle(failure);
			} catch (Exception e) {
			}
		}
	}
	
	@Override
	public Throwable getFailure(){
		return failure;
	}
	
	@Override
	public void setSussessHandler(Handler<T> handler){
		successHandler = handler;
		if(isCompleted)
			if(failure == null)
				setResult(result);
	}
	
	@Override
	public void setFailureHandler(Handler<Throwable> handler){
		failureHandler = handler;
		if(isCompleted)
			if(failure != null)
				setFailure(failure);
	}
	
	@Override
	public void setSuccess(){
		success = true;
	}
	
	@Override
	public void setFail(){
		success = false;
	}
	
	@Override
	public void setFail(Integer code, Throwable failure) {
		setFailure(failure);
		setErrorCode(code);
	}
	
	@Override
	public void setFail(Integer code, String message){
		setFail(code, new Throwable(message));
	}
	
	@Override
	public boolean isSuccess(){
		return success;
	}
	
	@Override
	public boolean isFail(){
		return !success;
	}

	@Override
	public void setErrorCode(Integer code) {
		this.errorCode = code;
	}

	@Override
	public Integer getErrorCode() {
		return errorCode;
	}

}
