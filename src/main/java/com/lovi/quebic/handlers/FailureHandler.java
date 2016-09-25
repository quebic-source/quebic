package com.lovi.quebic.handlers;

public interface FailureHandler<T>{
	void handle(T t, Throwable failure, int responseCode, String responseReason);
}
