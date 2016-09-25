package com.lovi.quebic.launcher;

import java.io.PrintStream;
import org.springframework.boot.Banner;
import org.springframework.boot.SpringApplication;
import org.springframework.context.ApplicationContext;
import org.springframework.core.env.Environment;

import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.servicecaller.MessengerMapperGenerator;
import com.lovi.quebic.servicecaller.impl.MessengerMapperGeneratorImpl;
import com.lovi.quebic.sockm.launcher.SockMLauncher;
import com.lovi.quebic.sockm.messenger.Messenger;
import com.lovi.quebic.web.HttpServer;
import com.lovi.quebic.web.HttpServerOption;
import com.lovi.quebic.web.RequestMapper;
import com.lovi.quebic.web.RequestMapperGenerator;
import com.lovi.quebic.web.impl.RequestMapperGeneratorImpl;

public class AppLauncherImpl implements AppLauncher {

	private static Object lock = new Object();
	private static AppLauncherImpl instance;
	private HttpServer httpServer;
	private SockMLauncher sockMLauncher;
	
	private final String DEFAULT_HOST_NAME = "localhost"; 
	
	private AppLauncherImpl() {
		httpServer = HttpServer.create();
	}
	
	public static AppLauncher create(){
		synchronized (lock) {
			if(instance == null)
				instance = new AppLauncherImpl();
		}
		return instance;
	}
	
	@Override
	public void addHttpServerOption(HttpServerOption option) {
		httpServer.addHttpServerOption(option);
	}
	
	@Override
	public void run(Class<?> baseClass, int port, String... args) {
		run(baseClass, DEFAULT_HOST_NAME, port, null, null, args);
	}

	@Override
	public void run(Class<?> baseClass, String hostname, int port,
			String... args) {
		run(baseClass, hostname, port, null, null, args);
	}

	@Override
	public void run(Class<?> baseClass, int port, Handler<String> successHandler,
			Handler<Throwable> failureHandler, String... args) {
		
		run(baseClass, DEFAULT_HOST_NAME, port, successHandler, failureHandler, args);
	}
	
	@Override
	public void run(Class<?> baseClass, String hostname, int port, 
			Handler<String> successHandler, Handler<Throwable> failureHandler,
			String... args) {
		Object[] objects = new Object[2];
		objects[0] = AppLauncher.class;
		objects[1] = baseClass;
		
		SpringApplication app = new SpringApplication(objects);
    	app.setBanner(new Banner() {
			
			@Override
			public void printBanner(Environment arg0, Class<?> arg1, PrintStream printStream) {
				printStream.println("##############");
				printStream.println("### quebic ###");
				printStream.println("##############");
			}
		});
        //app.setBannerMode(Banner.Mode.OFF);
    	ApplicationContext context = app.run(args);
    	

		/**
		 * Http Server Start
		 */
    	RequestMapper requestMapper = RequestMapper.create();

		RequestMapperGenerator requestMapperGenerator = new RequestMapperGeneratorImpl(context);
		try{
			requestMapperGenerator.start(baseClass, requestMapper);
		}catch(Exception e){

			if(failureHandler != null)
				failureHandler.handle(e);
			
			return;
		}
		
		startThread(()->{
			httpServer.run(hostname, port, requestMapper, successHandler, failureHandler);
		}, "quebic-http-server-container-thread");
		
		
		/**
		 * Http Server End
		 */
		
		/**
		 * SockM Start
		 */
		sockMLauncher = SockMLauncher.create();
		Messenger messenger = sockMLauncher.getMessenger();
		
		MessengerMapperGenerator messengerMapperGenerator = new MessengerMapperGeneratorImpl();
		try{
			messengerMapperGenerator.start(baseClass, messenger);
		}catch(Exception e){
			if(failureHandler != null)
				failureHandler.handle(e);
			
			return;
		}
		/**
		 * SockM End
		 */
		
	}
	
	private void startThread(Runnable runnable, String threadName){
		Thread thread = new Thread(runnable, threadName);
		thread.start();
	}
}
