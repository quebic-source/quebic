package com.lovi.quebic.cluster.message;

public class MessageKey {

	public static final String JOIN_NEW_SERVER_MEMBER = "JOIN_NEW_SERVER_MEMBER";
	public static final String WELCOME_NEW_SERVER_MEMBER = "WELCOME_NEW_SERVER_MEMBER";
	
	public static final String REQUEST_LIVE_DATA = "REQUEST_LIVE_DATA";
	public static final String LIVE_DATA = "LIVE_DATA";
	
	//session parm
	public static final String SESSION_NEW_USER = "SESSION_NEW_USER";
	public static final String SESSION_REMOVE_USER = "SESSION_REMOVE_USER";
	public static final String SESSION_PARM_PUT = "SESSION_PARM_PUT";
	public static final String SESSION_PARM_REMOVE = "SESSION_PARM_REMOVE";
	public static final String SESSION_PARMS_CLEAR_ALL = "SESSION_PARMS_CLEAR_ALL";
	
	//application parm
	public static final String APPLICATION_PARM_PUT = "APPLICATION_PARM_PUT";
	public static final String APPLICATION_PARM_REMOVE = "APPLICATION_PARM_REMOVE";
	public static final String APPLICATION_PARMS_CLEAR_ALL = "APPLICATION_PARMS_CLEAR_ALL";
	
	public static final String EXIT = "EXIT";
}
