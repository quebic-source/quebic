package com.lovi.quebic.web.server;

import com.lovi.quebic.web.HttpServer;

import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.socket.SocketChannel;
import io.netty.handler.codec.http.HttpContentCompressor;
import io.netty.handler.codec.http.HttpRequestDecoder;
import io.netty.handler.codec.http.HttpResponseEncoder;
import io.netty.handler.ssl.SslContext;

public class HttpServerInitializer extends ChannelInitializer<SocketChannel> {

    private final SslContext sslCtx;
    private HttpServer httpServer;

    public HttpServerInitializer(SslContext sslCtx, HttpServer httpServer) {
        this.sslCtx = sslCtx;
        this.httpServer = httpServer;
    }

    @Override
    public void initChannel(SocketChannel ch) {
    	ChannelPipeline pipeline = ch.pipeline();

        if (sslCtx != null) {
            pipeline.addLast(sslCtx.newHandler(ch.alloc()));
        }

        pipeline.addLast(new HttpRequestDecoder());
        pipeline.addLast(new HttpResponseEncoder());

        pipeline.addLast(new HttpContentCompressor());

        pipeline.addLast(new HttpServerHandler(httpServer));
    }
}