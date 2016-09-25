package com.lovi.quebic.cluster;

import java.net.InetSocketAddress;
import java.nio.channels.SocketChannel;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Executors;
import java.util.concurrent.ForkJoinPool;
import java.util.concurrent.RecursiveAction;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

class UnreachableMembersChecker extends RecursiveAction{
	
	private static final long serialVersionUID = 8687427634121686640L;
	
	final static Logger logger = LoggerFactory.getLogger(UnreachableMembersChecker.class);
	
	private static ClusterConnector groupConnector;
	
	private Member member;
	private boolean start;
	
	public UnreachableMembersChecker(boolean start) {
		this.start = start;
	}
	
	public UnreachableMembersChecker(boolean start, Member member) {
		this.start = start;
		this.member = member;
	}
	
	public static void startCheck(ClusterConnector groupConnector){
		
		Thread thread = new Thread(new Runnable() {
			
			@Override
			public void run() {
				UnreachableMembersChecker.groupConnector = groupConnector;
				
				ScheduledExecutorService scheduledExecutorService = Executors.newScheduledThreadPool(1);

				ScheduledFuture<?> scheduledFuture = scheduledExecutorService.scheduleWithFixedDelay(new Runnable() {
					
					@Override
					public void run() {
						
						UnreachableMembersChecker membersChecker = new UnreachableMembersChecker(true);
						
						ForkJoinPool forkJoinPool = new ForkJoinPool();
						forkJoinPool.execute(membersChecker);
						
						while (!membersChecker.isDone()){}
						
						forkJoinPool.shutdown();
						
					}
				}, 0, 1, TimeUnit.SECONDS);

				try {
					scheduledFuture.get();
				} catch (Exception e) {
					logger.error(e.getMessage());
				}
				
				scheduledExecutorService.shutdown();
			}
		});
		
		thread.start();
	}

	@Override
	protected void compute() {
		if(start)
			invokeAll(checkAction());
		else{
			
			try{
				InetSocketAddress socketAddress = new InetSocketAddress(member.getAddress(), member.getPort());
				SocketChannel socketChannel = SocketChannel.open(socketAddress);
				socketChannel.close();
				
			}catch(Exception e){
				groupConnector.removeTcpMemberByMember(member);
			}
			
			
		}
	}
	
	private List<UnreachableMembersChecker> checkAction(){
		
		List<UnreachableMembersChecker> actions = new ArrayList<>();
		
		for(Member member : groupConnector.getTcpMembers()){
			
			if(!member.equals(groupConnector.getLocalTcpMember())){
				UnreachableMembersChecker action = new UnreachableMembersChecker(false, member);
				actions.add(action);
			}
		}
		
		return actions;
	}
}
