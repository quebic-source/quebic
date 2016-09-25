package com.lovi.quebic.launcher;

import org.springframework.boot.autoconfigure.SpringBootApplication;

import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.web.HttpServerOption;

@SpringBootApplication
public interface AppLauncher {

	public static AppLauncher create() {
		return AppLauncherImpl.create();
	}

	void addHttpServerOption(HttpServerOption option);

	void run(Class<?> baseClass, int port, String... args);

	void run(Class<?> baseClass, String hostname, int port,
			String... args);

	void run(Class<?> baseClass, int port,
			Handler<String> successHandler, Handler<Throwable> failureHandler,
			String... args);

	void run(Class<?> baseClass, String hostname, int port, 
			Handler<String> successHandler, Handler<Throwable> failureHandler,
			String... args);

}
