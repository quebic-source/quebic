package com.lovi.quebic.common;

import java.util.concurrent.ThreadFactory;

public class ServerThreadFactory implements ThreadFactory {
	private int i;
	private String threadName;
	
	public ServerThreadFactory(String threadName) {
		this.i = 1;
		this.threadName = threadName;
	}
	
	@Override
	public Thread newThread(Runnable r) {
		return new Thread(r, threadName + "-" + (i++));
	}

}
