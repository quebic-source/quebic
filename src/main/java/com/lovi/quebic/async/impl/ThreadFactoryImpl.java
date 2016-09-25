package com.lovi.quebic.async.impl;

import java.util.concurrent.ThreadFactory;

public class ThreadFactoryImpl implements ThreadFactory{

	private int counter = 1;
	private String prefix;

	public ThreadFactoryImpl(String prefix) {
		this.prefix = prefix;
	}

	@Override
	public Thread newThread(Runnable r) {
		return new Thread(r, prefix + "-" + counter++);
	}
}
