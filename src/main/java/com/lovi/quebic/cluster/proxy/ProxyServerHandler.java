package com.lovi.quebic.cluster.proxy;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.cluster.Member;
import com.lovi.quebic.cluster.loadblancer.LoadBlancer;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import io.netty.util.ReferenceCountUtil;

public class ProxyServerHandler extends ChannelInboundHandlerAdapter{

	final static Logger logger = LoggerFactory.getLogger(ProxyServerHandler.class);
	
	private ByteBuf in;
	
	private LoadBlancer loadBlancer;
	
	public ProxyServerHandler(LoadBlancer loadBlancer){
		this.loadBlancer = loadBlancer; 
	}
 	
 	@Override
 	public void channelRead(ChannelHandlerContext ctx, Object msg) { 
       
 		if(in == null)
 			in = (ByteBuf) msg;
		try {
			
			
	        
		}catch(Exception e){ 
			logger.error("channelRead " + e.getMessage());
		}
		finally {
			ReferenceCountUtil.release(msg);
		}
 		
    }

 	@Override
 	public void channelReadComplete(ChannelHandlerContext ctx) throws Exception {
 		
 		try {
			
			Member destinationMember = loadBlancer.routeMember();//select destination member
			logger.info("select => \n" + destinationMember + "\n");
			
	        
		}catch(Exception e){ 
			logger.error("channelReadComplete " + e.getMessage());
		}
		finally {
		}
 		
 		
 	}
 	
    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) { // (4)
        logger.error(cause.getMessage());
        ctx.close();
    }
	
}
