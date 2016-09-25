package com.lovi.quebic.cluster.message;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.lovi.quebic.cluster.Member;

public class CoreMessage implements AutoCloseable{

	private String key;
	private Member member;
	private List<Member> members = new ArrayList<>();
	private Map<String, Map<String, Object>> sessionParm = new HashMap<>();
	private Map<String, Object> applicationContextParms = new HashMap<>();
	
	public String getKey() {
		return key;
	}
	public void setKey(String key) {
		this.key = key;
	}
	public Member getMember() {
		return member;
	}
	public void setMember(Member member) {
		this.member = member;
	}
	public List<Member> getMembers() {
		return members;
	}
	public void setMembers(List<Member> members) {
		this.members = members;
	}
	public Map<String, Map<String, Object>> getSessionParm() {
		return sessionParm;
	}
	public void setSessionParm(Map<String, Map<String, Object>> sessionParm) {
		this.sessionParm = sessionParm;
	}
	public Map<String, Object> getApplicationContextParms() {
		return applicationContextParms;
	}
	public void setApplicationContextParms(
			Map<String, Object> applicationContextParms) {
		this.applicationContextParms = applicationContextParms;
	}
	@Override
	public void close() throws Exception {
	}
	
}
