package com.lovi.quebic.async;

import com.lovi.quebic.async.impl.FutureImpl;
import com.lovi.quebic.handlers.Handler;

/**
 * 
 * @author Tharanga Thennakoon
 *
 * @param <R>
 */
public interface Future<T> {
	
	static <T> Future<T> create(){
		return new FutureImpl<>();
	}
	
    void setResult(T result);
    
	T getResult();
 
    void setFailure(Throwable failure);

	Throwable getFailure();
    
	void setSussessHandler(Handler<T> handler);

	void setFailureHandler(Handler<Throwable> handler);

	void setSuccess();

	void setFail();
	
	void setFail(Integer code, Throwable failure);
	
	void setFail(Integer code, String message);
	
	boolean isSuccess();

	boolean isFail();

	void setErrorCode(Integer code);
	
	Integer getErrorCode();
	
}
