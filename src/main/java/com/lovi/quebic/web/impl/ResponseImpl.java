package com.lovi.quebic.web.impl;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelFutureListener;
import io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.netty.handler.codec.http.FullHttpResponse;
import io.netty.handler.codec.http.HttpHeaderNames;
import io.netty.handler.codec.http.HttpHeaderValues;
import io.netty.handler.codec.http.HttpRequest;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.codec.http.HttpVersion;
import java.util.Date;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Map.Entry;
import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.exception.JsonParserException;
import com.lovi.quebic.web.Response;

public class ResponseImpl implements Response{
	
	final static Logger logger = LoggerFactory.getLogger(ResponseImpl.class);

	private Channel channel;
	private HttpRequest httpRequest;
	
    private int responseCode = 200;
    private String responseReason = "OK";
    private Map<String, String> headers = new LinkedHashMap<String, String>();
    private byte[] content;
    
    /**
     * Use chunkedTransfer
     */
    private boolean chunkedTransfer;
    private boolean encodeAsGzip;
    private boolean keepAlive;
   
    public ResponseImpl(Channel channel, HttpRequest httpRequest) {
    	
    	this.channel = channel;
    	this.httpRequest = httpRequest;
    	
        headers.put("Date", new Date().toString());
        headers.put("Server", "quebic web server");
    }

    @Override
    public int getResponseCode() {
        return responseCode;
    }

    @Override
    public String getResponseReason() {
        return responseReason;
    }
    
    @Override
    public void setResponseCode(int responseCode) {
        this.responseCode = responseCode;
    }

    @Override
    public void setResponseReason(String responseReason) {
        this.responseReason = responseReason;
    }

    @Override
    public String getHeader(String header) {
        return headers.get(header);
    }
    
    @Override
    public void setHeader(String key, String value) {
        headers.put(key, value);
    }

    @Override
    public Map<String, String> getHeaders() {
        return headers;
    }
    
    @Override
    public void setContentType(String contentType){
    	setHeader("Content-Type", contentType);
    }
    
    @Override
    public String getContentType(){
    	return getHeader("Content-Type");
    }
    
    @Override
    public byte[] getContent() {
        return content;
    }

    @Override
    public void setEncodeAsGzip(boolean encodeAsGzip) {
		this.encodeAsGzip = encodeAsGzip;
	}

    @Override
	public void setKeepAlive(boolean keepAlive) {
		this.keepAlive = keepAlive;
	}
    
    @Override
	public void setChunkedTransfer(boolean chunkedTransfer) {
		this.chunkedTransfer = chunkedTransfer;
	}
    
    @Override
	public boolean isChunkedTransfer() {
		return chunkedTransfer;
	}

    @Override
	public boolean isEncodeAsGzip() {
		return encodeAsGzip;
	}

    @Override
	public boolean isKeepAlive() {
		return keepAlive;
	}

	@Override
    public void writePOJO(Object content) throws JsonParserException{
    	
    	try{
    		ObjectMapper objectMapper = new ObjectMapper();
    		String responseValue = objectMapper.writeValueAsString(content);
    		
    		setContentType("application/json");
            write(responseValue.getBytes());
    	}catch(Exception e){
    		throw new JsonParserException(e);
    	}
    	
    }
    
    @Override
    public void write(Integer content) {
    	write(String.valueOf(content));
    }
    
    @Override
    public void write(Short content) {
    	write(String.valueOf(content));
    }
    
    @Override
    public void write(Long content) {
    	write(String.valueOf(content));
    }
    
    @Override
    public void write(Float content) {
    	write(String.valueOf(content));
    }
    
    @Override
    public void write(Double content) {
    	write(String.valueOf(content));
    }
    
    @Override
    public void write(Boolean content) {
    	write(String.valueOf(content));
    }
    
    @Override
    public void write(String content) {
    	write(content.getBytes());
    }
    
    @Override
    public void write(byte[] content) {
    	
    	try {
        	
            ByteBuf buf = Unpooled.wrappedBuffer(content);
            
            // Decide whether to close the connection or not.
            boolean close = httpRequest.headers().contains(HttpHeaderNames.CONNECTION, HttpHeaderValues.CLOSE, true)
                    || httpRequest.protocolVersion().equals(HttpVersion.HTTP_1_0)
                    && !httpRequest.headers().contains(HttpHeaderNames.CONNECTION, HttpHeaderValues.KEEP_ALIVE, true);

            // Build the response object.
            FullHttpResponse response = new DefaultFullHttpResponse(
                    HttpVersion.HTTP_1_1, new HttpResponseStatus(responseCode, responseReason), buf);
            
            for(Entry<String, String> entry : headers.entrySet()){
            	response.headers().add(entry.getKey(), entry.getValue());
            }
            
            if (!close) {
                // There's no need to add 'Content-Length' header
                // if this is the last response.
                response.headers().setInt(HttpHeaderNames.CONTENT_LENGTH, buf.readableBytes());
            }
            
            // Write the response.
            ChannelFuture future = channel.writeAndFlush(response);
            // Close the connection after the write operation is done if necessary.
            if (close) {
                future.addListener(ChannelFutureListener.CLOSE);
            }
           
        } catch (Exception ex) {
        	logger.error(ex.getMessage());
        }
    	
    }

}
