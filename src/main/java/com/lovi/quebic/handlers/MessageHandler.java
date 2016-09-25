package com.lovi.quebic.handlers;

import com.lovi.quebic.async.Future;

public interface MessageHandler<T>{
	void handle(T t, Future<?> future);
}
