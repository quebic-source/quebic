package com.lovi.quebic.cluster.loadblancer;

import java.util.List;

import com.lovi.quebic.cluster.Member;
import com.lovi.quebic.cluster.loadblancer.impl.RoundRobinLoadBlancer;
import com.lovi.quebic.exception.WorkerServesNotFoundException;

public interface LoadBlancer {
	
	static LoadBlancer createDefaultLoadBlancer() {
		return new RoundRobinLoadBlancer();
	}
	
	void setMembers(List<Member> members);
	Member routeMember() throws WorkerServesNotFoundException;
	void resetRoute();
}
