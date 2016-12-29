package com.lovi.quebic.sockm.messenger.impl;

import java.util.HashMap;
import java.util.Map;
import java.util.Map.Entry;

import com.lovi.quebic.sockm.message.messenger.Message;
import com.lovi.quebic.sockm.messenger.Publisher;
import com.lovi.quebic.sockm.messenger.Subscriber;

public class PublisherImpl implements Publisher {

	private Map<String, Subscriber<?,?>> subscribersMap = new HashMap<>();
	
	public PublisherImpl() {
	}
	
	@Override
	public <I,O> void addSubscriber(String subscriberId,Subscriber<I,O> subscriber){
		subscribersMap.put(subscriberId, subscriber);
	}
	
	@Override
	public <I> void publish(String address, I message){
		for(Entry<String, Subscriber<?,?>> subscriberEntry : subscribersMap.entrySet()){
			
			Subscriber<?,?> subscriber = subscriberEntry.getValue();
			
			if(subscriber.getAddress().equals(address))
				subscriber.run(message);
			
		}
	}
	
	@SuppressWarnings("unchecked")
	@Override
	public <I,O> boolean send(String address, Message<I,O> message){
		
		boolean isFound = false;
		
		for(Entry<String, Subscriber<?,?>> subscriberEntry : subscribersMap.entrySet()){
			
			Subscriber<I,O> subscriber = (Subscriber<I, O>) subscriberEntry.getValue();
			
			if(subscriber.getAddress().equals(address)){
				isFound = true;
				subscriber.run(message);
				break;
			}
			
		}
		
		return isFound;
	}
	
	@SuppressWarnings("unchecked")
	@Override
	public <I,O> boolean sendBySubscriberId(String subscriberId, Message<I,O> message){
		Subscriber<I,O> subscriber = (Subscriber<I, O>) subscribersMap.get(subscriberId);
		if(subscriber != null){
			subscriber.run(message);
			return true;
		}else
			return false;
	}
}
