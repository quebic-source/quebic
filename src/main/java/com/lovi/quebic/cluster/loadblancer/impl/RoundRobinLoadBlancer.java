package com.lovi.quebic.cluster.loadblancer.impl;

import java.util.ArrayList;
import java.util.List;

import com.lovi.quebic.cluster.Member;
import com.lovi.quebic.cluster.loadblancer.LoadBlancer;
import com.lovi.quebic.exception.WorkerServesNotFoundException;

public class RoundRobinLoadBlancer implements LoadBlancer {

	private int selectedMember = 0;
	private List<Member> members = new ArrayList<>();
	
	@Override
	public void setMembers(List<Member> members) {
		this.members = members;
	}

	@Override
	public Member routeMember() throws WorkerServesNotFoundException {
		
		if(members.size() == 0)
			throw new WorkerServesNotFoundException();
		
		if(selectedMember >= members.size()){
			selectedMember = 0;
		}
		return members.get(selectedMember++);
	}

	@Override
	public void resetRoute(){
		selectedMember = 0;
	}
}
