package com.lovi.quebic.cluster.proxy;

import com.lovi.quebic.cluster.Member;

import io.netty.bootstrap.Bootstrap;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioSocketChannel;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslContextBuilder;
import io.netty.handler.ssl.util.InsecureTrustManagerFactory;

public class ProxyClient {

	static final boolean SSL = System.getProperty("ssl") != null;
	
	private Member member;
	
	public ProxyClient(Member member){
		this.member = member;
	}
	
	public void start() throws Exception{
		
		// Configure SSL.git
        final SslContext sslCtx;
        if (SSL) {
            sslCtx = SslContextBuilder.forClient()
                .trustManager(InsecureTrustManagerFactory.INSTANCE).build();
        } else {
            sslCtx = null;
        }

        EventLoopGroup group = new NioEventLoopGroup();
        try {
            Bootstrap b = new Bootstrap();
            b.group(group)
             .channel(NioSocketChannel.class)
             .option(ChannelOption.TCP_NODELAY, true)
             .handler(new ChannelInitializer<SocketChannel>() {
                 @Override
                 public void initChannel(SocketChannel ch) throws Exception {
                     
                	 ChannelPipeline p = ch.pipeline();
                     if (sslCtx != null) {
                         p.addLast(sslCtx.newHandler(ch.alloc(), member.getAddress(), member.getPort()));
                     }
                     
                     p.addLast(new ChannelInboundHandlerAdapter(){
                    	 
                    	 	private ByteBuf byteBufOut;
                    	
                    	 	@Override
                    	    public void channelActive(ChannelHandlerContext ctx) {
                    	 		
                    	 		
                    	    }

                    	    @Override
                    	    public void channelRead(ChannelHandlerContext ctx, Object msg) {
                    	       
                    	    }

                    	    @Override
                    	    public void channelReadComplete(ChannelHandlerContext ctx) {
                    	    	
                    	    }

                    	    @Override
                    	    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) {
                    	    	
                    	    }
                    	 
                    	 
                     });
                 }
             });

           
            ChannelFuture f = b.connect(member.getAddress(), member.getPort()).sync();

            
            f.channel().closeFuture().sync();
        } finally {
            group.shutdownGracefully();
        }
		
	}
	
}
