package com.lovi.quebic.handlers;

import com.lovi.quebic.async.Future;

public interface RequestHandler<T>{
	void handle(T t, Future<?> future);
}
