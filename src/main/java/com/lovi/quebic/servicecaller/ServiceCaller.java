package com.lovi.quebic.servicecaller;

import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.servicecaller.impl.ServiceCallerImpl;

public interface ServiceCaller<I,O>{

	public static <I,O> ServiceCaller<I,O> create(String address) {
		return new ServiceCallerImpl<I,O>(address);
	}
	
	public static <I,O> ServiceCaller<I,O> create(String address, Class<?> inputClassType, Class<?> outputClassType) {
		return new ServiceCallerImpl<I,O>(address, inputClassType, outputClassType);
	}

	ServiceCaller<I, O> result(Handler<O> resultHandler);

	ServiceCaller<I, O> failure(Handler<Throwable> failureHandler);
	
	void call(Object... objects);


	public default ServiceCaller<O,?> andThen(ServiceCaller<O,?> nextCaller){
		result(r->{
			nextCaller.call(r);
		});
		return nextCaller;
	}

}
