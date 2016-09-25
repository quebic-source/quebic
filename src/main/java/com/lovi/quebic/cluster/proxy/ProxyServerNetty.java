package com.lovi.quebic.cluster.proxy;

import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslContextBuilder;
import io.netty.handler.ssl.util.InsecureTrustManagerFactory;
import io.netty.handler.ssl.util.SelfSignedCertificate;

import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.channels.ServerSocketChannel;
import java.util.Iterator;
import java.util.Set;
import java.util.concurrent.CompletableFuture;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.async.Future;
import com.lovi.quebic.cluster.Member;
import com.lovi.quebic.cluster.loadblancer.LoadBlancer;
import com.lovi.quebic.common.Message;

import io.netty.bootstrap.Bootstrap;
import io.netty.bootstrap.ServerBootstrap;
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
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.channel.socket.nio.NioSocketChannel;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import io.netty.util.ReferenceCountUtil;

public class ProxyServerNetty {

	final static Logger logger = LoggerFactory.getLogger(ProxyServerNetty.class);
	
	static final boolean SSL = System.getProperty("ssl") != null;
	private String sourceAddress;
	private int sourcePort;
	
	private LoadBlancer loadBlancer;
	
	public ProxyServerNetty(String sourceAddress,int sourcePort, LoadBlancer loadBlancer) {
		this.sourceAddress = sourceAddress;
		this.sourcePort = sourcePort;
		this.loadBlancer = loadBlancer;
	}
	
	public void start(Future<String> future) throws Exception {
		// Configure SSL.
        final SslContext sslCtx;
        if (SSL) {
            SelfSignedCertificate ssc = new SelfSignedCertificate();
            sslCtx = SslContextBuilder.forServer(ssc.certificate(), ssc.privateKey()).build();
        } else {
            sslCtx = null;
        }

        EventLoopGroup bossGroup = new NioEventLoopGroup(1);
        EventLoopGroup workerGroup = new NioEventLoopGroup();
        try {
            ServerBootstrap b = new ServerBootstrap();
            b.group(bossGroup, workerGroup)
             .channel(NioServerSocketChannel.class)
             .handler(new LoggingHandler(LogLevel.INFO))
             .childHandler(new ChannelInitializer<SocketChannel>() {
                 @Override
                 public void initChannel(SocketChannel ch) {
                     ChannelPipeline p = ch.pipeline();
                     if (sslCtx != null) {
                         p.addLast(sslCtx.newHandler(ch.alloc()));
                     }
                     p.addLast(new ChannelInboundHandlerAdapter(){
                    	 	private ByteBuf in;
                    	 	
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
                    				startClient(destinationMember, ctx, in);
                    				
                    				
                    		        
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
                    	 
                     });
                 }
             });

            ChannelFuture f = b.bind(sourcePort).sync();

            f.channel().closeFuture().sync();
        } finally {
            workerGroup.shutdownGracefully();
            bossGroup.shutdownGracefully();
        }
		
	}

	
	private void startClient(Member member, ChannelHandlerContext ctxIn, ByteBuf byteBufIn) throws Exception{

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
