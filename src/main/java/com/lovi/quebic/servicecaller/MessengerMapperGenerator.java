package com.lovi.quebic.servicecaller;

import com.lovi.quebic.sockm.messenger.Messenger;

public interface MessengerMapperGenerator {

	void start(Class<?> baseClass, Messenger messenger) throws Exception;

}
