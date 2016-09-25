package com.lovi.quebic.cluster.proxy;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;

public class ProxyClientHandler extends ChannelInboundHandlerAdapter{

	final static Logger logger = LoggerFactory.getLogger(ProxyClientHandler.class);
	
	private ByteBuf byteBufIn;
	private ChannelHandlerContext ctxIn;
	
	private ByteBuf byteBufOut;
	
	public ProxyClientHandler(ByteBuf byteBufIn, ChannelHandlerContext ctxIn){
		this.byteBufIn = byteBufIn;
		this.ctxIn = ctxIn;
	}
	
 	@Override
    public void channelActive(ChannelHandlerContext ctx) {
 		
 		byteBufOut = Unpooled.buffer(byteBufIn.capacity());
 		
 		while (byteBufIn.isReadable()) {
			byteBufOut.writeByte(byteBufIn.readByte());
		}
 		
 		ctx.writeAndFlush(byteBufOut);
    }

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) {
        ctxIn.write(msg);
    }

    @Override
    public void channelReadComplete(ChannelHandlerContext ctx) {
    	ctxIn.flush();
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) {
    	logger.error(cause.getMessage());
    	ctxIn.close();
        ctx.close();
    }
	
}
