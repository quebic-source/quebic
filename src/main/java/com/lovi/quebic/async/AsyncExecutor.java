package com.lovi.quebic.async;

import com.lovi.quebic.async.impl.AsyncExecutorImpl;
import com.lovi.quebic.handlers.AsyncHandler;
import com.lovi.quebic.handlers.Handler;

/**
 * 
 * @author Tharanga Thennakoon
 *
 */
public interface AsyncExecutor<T>{
	
	static <T> AsyncExecutor<T> create(){
		return new AsyncExecutorImpl<T>();
	}
	void run(AsyncHandler<T> handler, Handler<T> successHandler, Handler<Throwable> failureHandler);
}
