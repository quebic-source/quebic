package com.lovi.quebic.common;

public enum Message {

	SERVER_START("quebic webserver listen at : "),
	SERVER_CLUSTER_START("quebic webserver clustering starting..."),
	MASTER_SERVER_CLUSTER_START("quebic master webserver clustering starting..."),
	MASTER_SERVER_START("quebic master webserver listen "),
	SERVER_TERMINATE("SERVER_TERMINATE"),
	MEMBER_LEAVE("MEMBER_LEAVE "),
	UNREACHABLE_CHECK("UNREACHABLE_CHECK"),
	CLUSTER_FIRST_MEMBER("I am the first member"),
	CLUSTER_FIRST_WORKER_MEMBER("I am the first worker member. No any live data"),
	CLUSTER_MEMBER_WAITING_JOIN("waiting 2s for joining ..."),
	CLUSTER_WORKER_MEMBER_WAITING_LIVE_DATA("waiting 2s for getting live data...");
	
	private String message;
	private Message(String message){
		this.message = message;
	}
	
	public String getMessage(){
		return message;
	}
}
