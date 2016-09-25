package com.lovi.quebic.web.server;

import io.netty.buffer.ByteBuf;
import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.netty.handler.codec.http.FullHttpResponse;
import io.netty.handler.codec.http.HttpContent;
import io.netty.handler.codec.http.HttpHeaderNames;
import io.netty.handler.codec.http.HttpUtil;
import io.netty.handler.codec.http.HttpHeaderValues;
import io.netty.handler.codec.http.HttpMethod;
import io.netty.handler.codec.http.HttpObject;
import io.netty.handler.codec.http.HttpRequest;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.codec.http.HttpVersion;
import io.netty.handler.codec.http.LastHttpContent;
import io.netty.handler.codec.http.QueryStringDecoder;
import io.netty.handler.codec.http.multipart.Attribute;
import io.netty.handler.codec.http.multipart.DefaultHttpDataFactory;
import io.netty.handler.codec.http.multipart.DiskAttribute;
import io.netty.handler.codec.http.multipart.DiskFileUpload;
import io.netty.handler.codec.http.multipart.FileUpload;
import io.netty.handler.codec.http.multipart.HttpData;
import io.netty.handler.codec.http.multipart.HttpDataFactory;
import io.netty.handler.codec.http.multipart.HttpPostRequestDecoder;
import io.netty.handler.codec.http.multipart.HttpPostRequestDecoder.EndOfDataDecoderException;
import io.netty.handler.codec.http.multipart.HttpPostRequestDecoder.ErrorDataDecoderException;
import io.netty.handler.codec.http.multipart.InterfaceHttpData;
import io.netty.handler.codec.http.multipart.InterfaceHttpData.HttpDataType;
import io.netty.util.CharsetUtil;

import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.Map.Entry;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import com.lovi.quebic.async.Future;
import com.lovi.quebic.cluster.ClusterConnector;
import com.lovi.quebic.exception.ErrorMessage;
import com.lovi.quebic.exception.RequestMapperException;
import com.lovi.quebic.exception.ResourceNotFoundException;
import com.lovi.quebic.handlers.FailureHandler;
import com.lovi.quebic.web.ApplicationContextData;
import com.lovi.quebic.web.HttpServer;
import com.lovi.quebic.web.Request;
import com.lovi.quebic.web.RequestMap;
import com.lovi.quebic.web.RequestMapper;
import com.lovi.quebic.web.Response;
import com.lovi.quebic.web.ServerContext;
import com.lovi.quebic.web.Session;
import com.lovi.quebic.web.SessionStore;
import com.lovi.quebic.web.file.HttpFile;
import com.lovi.quebic.web.file.impl.HttpFileAttributeImpl;
import com.lovi.quebic.web.impl.ClusterSessionImpl;
import com.lovi.quebic.web.impl.RequestImpl;
import com.lovi.quebic.web.impl.ResponseImpl;
import com.lovi.quebic.web.impl.ServerContextImpl;
import com.lovi.quebic.web.impl.SessionImpl;

import static io.netty.buffer.Unpooled.*;

public class HttpServerHandler extends SimpleChannelInboundHandler<HttpObject> {

    private static final Logger logger = Logger.getLogger(HttpServerHandler.class.getName());

    private HttpRequest httpRequest;

    private boolean readingChunks;

    private HttpData partialContent;

    private final StringBuilder responseContent = new StringBuilder();

    private static final HttpDataFactory factory =
            new DefaultHttpDataFactory(DefaultHttpDataFactory.MINSIZE); // Disk if size exceed

    private HttpPostRequestDecoder decoder;
    
    //quebic
    private HttpServer httpServer;
    private RequestMapper requestMapper;
    private ClusterConnector clusterConnector;
	private SessionStore sessionStore;
	private ApplicationContextData applicationContextData;
	
	private RequestImpl requestNettyImpl;
	//quebic

    static {
        DiskFileUpload.deleteOnExitTemporaryFile = true; // should delete file
                                                         // on exit (in normal
                                                         // exit)
        DiskFileUpload.baseDirectory = null; // system temp directory
        DiskAttribute.deleteOnExitTemporaryFile = true; // should delete file on
                                                        // exit (in normal exit)
        DiskAttribute.baseDirectory = null; // system temp directory
    }

    public HttpServerHandler(HttpServer httpServer) {
    	
		this.httpServer = httpServer;
		
		this.requestMapper = httpServer.getRequestMapper();
		this.clusterConnector = httpServer.getClusterConnector();
		
		this.applicationContextData = this.httpServer.getApplicationContextData();
		this.sessionStore = this.httpServer.getSessionStore();
		
	}
    
    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
        if (decoder != null) {
            decoder.cleanFiles();
        }
    }

    @Override
    public void channelRead0(ChannelHandlerContext ctx, HttpObject msg) throws Exception {
    	if (msg instanceof HttpRequest) {
    		HttpRequest httpRequest = this.httpRequest = (HttpRequest) msg;
    		
    		if(requestNettyImpl == null)
    			requestNettyImpl = new RequestImpl();
    		
    		//location
    		requestNettyImpl.setLocation(httpRequest.uri());
    		
    		//headers
    		for(Entry<String, String> entry : httpRequest.headers()) {
    			requestNettyImpl.getHeaders().put(entry.getKey(), entry.getValue());
            }
    		
    		//query parms
    		QueryStringDecoder decoderQuery = new QueryStringDecoder(httpRequest.uri());
            Map<String, List<String>> uriAttributes = decoderQuery.parameters();
            for(Entry<String, List<String>> attr: uriAttributes.entrySet()) {
            	requestNettyImpl.getParameters().put(attr.getKey(), attr.getValue().get(0));
            }
            
            //http method
    		requestNettyImpl.setMethodStr(httpRequest.method().name());
    		
    		// if GET Method: should not try to create a HttpPostRequestDecoder
            if (httpRequest.method().equals(HttpMethod.GET)) {
            	requestProcess(ctx.channel());
                return;
            }
            try {
                decoder = new HttpPostRequestDecoder(factory, httpRequest);
            } catch (ErrorDataDecoderException ex) {
                writeServerErrorResponse(ctx.channel() , ex);
                ctx.channel().close();
                return;
            }
            
            readingChunks = HttpUtil.isTransferEncodingChunked(httpRequest);
            if (readingChunks) {
                readingChunks = true;
            }
    		
    	}
    	
        if (decoder != null) {
            if (msg instanceof HttpContent) {
            	
                HttpContent chunk = (HttpContent) msg;
                try {
                    decoder.offer(chunk);
                } catch (ErrorDataDecoderException e1) {
                    e1.printStackTrace();
                    writeServerErrorResponse(ctx.channel() , e1);
                    ctx.channel().close();
                    return;
                }
                
                readHttpDataChunkByChunk();
                
                if (chunk instanceof LastHttpContent) {
                    
                    readingChunks = false;
                    reset();
                    requestProcess(ctx.channel());
                }
            }
        }
    }
    
    private void reset() {
        decoder.destroy();
        decoder = null;
    }

    
    private void readHttpDataChunkByChunk() {
        try {
            while (decoder.hasNext()) {
                InterfaceHttpData data = decoder.next();
                if (data != null) {
                    // check if current HttpData is a FileUpload and previously set as partial
                    if (partialContent == data) {
                        //logger.info(" 100% (FinalSize: " + partialContent.length() + ")");
                        partialContent = null;
                    }
                    try {
                       
                        writeHttpData(data);
                    } finally {
                        data.release();
                    }
                }
            }
            
            InterfaceHttpData data = decoder.currentPartialHttpData();
            if (data != null) {
                StringBuilder builder = new StringBuilder();
                if (partialContent == null) {
                    partialContent = (HttpData) data;
                    if (partialContent instanceof FileUpload) {
                        builder.append("Start FileUpload: ")
                            .append(((FileUpload) partialContent).getFilename()).append(" ");
                    } else {
                        builder.append("Start Attribute: ")
                            .append(partialContent.getName()).append(" ");
                    }
                    builder.append("(DefinedSize: ").append(partialContent.definedLength()).append(")");
                }
                if (partialContent.definedLength() > 0) {
                    builder.append(" ").append(partialContent.length() * 100 / partialContent.definedLength())
                        .append("% ");
                    //logger.info(builder.toString());
                } else {
                    builder.append(" ").append(partialContent.length()).append(" ");
                    //logger.info(builder.toString());
                }
            }
        } catch (EndOfDataDecoderException e1) {
        	//System.err.println("readHttpDataChunkByChunk " + e1.getMessage());
        }
    }

    private void writeHttpData(InterfaceHttpData data) {
        if (data.getHttpDataType() == HttpDataType.Attribute) {
            Attribute attribute = (Attribute) data;
            String value;
            try {
                value = attribute.getValue();
            } catch (IOException e1) {
                e1.printStackTrace();
                return;
            }
            
            if (value.length() > 100) {
            	logger.info("data too long " + attribute);
            }
            
            requestNettyImpl.getParameters().put(attribute.getName(), value);
            
            
        } else {
        	
            if (data.getHttpDataType() == HttpDataType.FileUpload) {
                FileUpload fileUpload = (FileUpload) data;
                if (fileUpload.isCompleted()) {
                    
                	/*if (fileUpload.length() < 10000) {
                        
                        try {
                        	System.out.println("getCharset()");
                        	System.out.println(fileUpload.getString(fileUpload.getCharset()));
                        } catch (IOException e1) {
                            e1.printStackTrace();
                        }
                    } else {
                    	System.out.println("File too long to be printed out : " + fileUpload.length());
                    }*/
                	
                    // fileUpload.isInMemory();// tells if the file is in Memory
                    // or on File
                    try{
                    	//fileUpload.renameTo(new File(fileUpload.getFilename())); // enable to move into another
                    	
                    	HttpFile httpFile = new HttpFileAttributeImpl(fileUpload);
                    	
                    	requestNettyImpl.getFiles().put(fileUpload.getName(), httpFile);
                    }catch(Exception e){
                    	System.err.println("file upload " + e.getMessage());
                    }
                    // File dest
                    // decoder.removeFileUploadFromClean(fileUpload); //remove
                    // the File of to delete file
                } else {
                	System.out.println("File to be continued but should not!");
                }
            }
        }
    }

   

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        logger.log(Level.WARNING, responseContent.toString(), cause);
        ctx.channel().close();
    }
    
    private void requestProcess(Channel channel) throws Exception {
		
		Request requst = this.requestNettyImpl;
		Response httpResponse = new ResponseImpl(channel, httpRequest);
		Session session = createSession(requst, httpResponse);

		// -------------------------mapping request----------------
		ServerContext serverContext = new ServerContextImpl(httpServer, this, requst, httpResponse, applicationContextData, session);

		if (requestMapper != null) {
			String userRequestPath = null;
			boolean found = false;
			for (RequestMap requestMap : requestMapper.getRequestMaps()) {

				String requestPath = requestMap.getPath();
				String regExpPath = requestMap.getRegExpPath();

				String incomePath = requst.getLocation().split("\\?")[0];// remove
																				// query
																				// parameters;
				incomePath = (incomePath.charAt(incomePath.length() - 1) == '/')
						? incomePath.substring(0, incomePath.length() - 1) : incomePath;// remove last slash

						
				userRequestPath = incomePath;
						
				// add index point
				requestPath = "/_INDEX_" + requestPath;
				regExpPath = "/_INDEX_" + regExpPath;
				incomePath = "/_INDEX_" + incomePath;

				String[] splitRequestPath = requestPath.split("/");
				String[] splitImcomePathStr = incomePath.split("/");

				// http method are not equal
				if (!requestMap.getHttpMethod().toString().equals(requst.getMethodStr()))
					continue;

				// check route and incoming request path separators not equal and this is not wildcard path
				if ((splitRequestPath.length != splitImcomePathStr.length) && !splitRequestPath[splitRequestPath.length - 1].equals("*"))
					continue;

				// match
				Matcher matcher = Pattern.compile(regExpPath).matcher(incomePath);
				if (matcher.matches()) {

					// extract path variables
					for (int i = 0; i < splitRequestPath.length; i++) {

						String str = splitRequestPath[i];

						Matcher matcherParm = Pattern.compile("\\{(.*)\\}").matcher(str);

						if (matcherParm.matches()) {

							try {
								serverContext.getHttpRequst().getParameters().put(matcherParm.group(1),
										splitImcomePathStr[i]);
							} catch (Exception e) {
								continue;
							}
						}

					}

					//handler process
					Future<?> future = Future.create();
					try {
						requestMap.getHandler().handle(serverContext, future);
						
						if(future.isFail())
							throw future.getFailure();
						
					} catch (Throwable e) {
						
						int responseCode = 500;
						if(future.getErrorCode() != null)
							responseCode = future.getErrorCode();
						
						FailureHandler<ServerContext> failureHandlerRequestMap = requestMap.getFailureHandler();

						if (failureHandlerRequestMap != null)
							failureHandlerRequestMap.handle(serverContext, e, responseCode,
									ErrorMessage.INTERNAL_SERVER_ERROR.getMessage());
						else {
							FailureHandler<ServerContext> failureHandlerRequestMapper = requestMapper
									.getFailureHandler();
							if (failureHandlerRequestMapper != null)
								failureHandlerRequestMapper.handle(serverContext, e, responseCode,
										ErrorMessage.INTERNAL_SERVER_ERROR.getMessage());
							else
								writeServerError(serverContext, e, responseCode,
										ErrorMessage.INTERNAL_SERVER_ERROR.getMessage());
						}
					}
					found = true;
					break;
					
				}
			}

			if (!found) {
				FailureHandler<ServerContext> failureHandlerRequestMapper = requestMapper.getFailureHandler();
				if (failureHandlerRequestMapper != null)
					failureHandlerRequestMapper.handle(
							serverContext, 
							new ResourceNotFoundException(ErrorMessage.RESOURCE_NOT_FOUND.getMessage() + " " + userRequestPath), 
							404,
							ErrorMessage.RESOURCE_NOT_FOUND.getMessage());
				else
					writeServerError(
							serverContext, 
							new ResourceNotFoundException(ErrorMessage.RESOURCE_NOT_FOUND.getMessage() + " " + userRequestPath), 
							404,
							ErrorMessage.RESOURCE_NOT_FOUND.getMessage());
			}

		} else {
			throw new RequestMapperException();
		}

	}

    /**
	 * failure -> RequestMap [if not handle] -> RequestMapper -> [if not handle]
	 * -> writeServerError()
	 * 
	 * @param serverContext
	 * @param failure
	 * @param responseCode
	 * @param responseReason
	 */
	private void writeServerError(ServerContext serverContext, Throwable failure, int responseCode,
			String responseReason) {
		Response response = serverContext.getHttpResponse();
		response.setResponseCode(responseCode);
		response.setResponseReason(responseReason);

		StringBuilder responseStrBuilder = new StringBuilder();
		responseStrBuilder
				.append("<h1 style='background-color:#D50000;color:#FFF'>HTTP Status - " + responseCode + "</h1>");
		responseStrBuilder.append("<h3 style='color:#D50000'>message : " + failure.getMessage() + "</h3>");
		responseStrBuilder.append("<h3 style='color:#D50000'>puppy-io [web]</h3>");

		response.setHeader("Content-type", "text/html");
		response.write(responseStrBuilder.toString());
	}
	
	 private void writeServerErrorResponse(Channel channel, Throwable failure) {
	    	
	    	
	    	StringBuilder responseStrBuilder = new StringBuilder();
			responseStrBuilder
					.append("<h1 style='background-color:#D50000;color:#FFF'>HTTP Status - " + 500 + "</h1>");
			responseStrBuilder.append("<h3 style='color:#D50000'>message : " + failure.getMessage() + "</h3>");
			responseStrBuilder.append("<h3 style='color:#D50000'>puppy-io [web]</h3>");
	    	
	        // Convert the response content to a ChannelBuffer.
	        ByteBuf buf = copiedBuffer(responseContent.toString(), CharsetUtil.UTF_8);
	        responseContent.setLength(0);

	        // Decide whether to close the connection or not.
	        boolean close = httpRequest.headers().contains(HttpHeaderNames.CONNECTION, HttpHeaderValues.CLOSE, true)
	                || httpRequest.protocolVersion().equals(HttpVersion.HTTP_1_0)
	                && !httpRequest.headers().contains(HttpHeaderNames.CONNECTION, HttpHeaderValues.KEEP_ALIVE, true);

	        // Build the response object.
	        FullHttpResponse response = new DefaultFullHttpResponse(
	                HttpVersion.HTTP_1_1, HttpResponseStatus.INTERNAL_SERVER_ERROR, buf);
	        response.headers().set(HttpHeaderNames.CONTENT_TYPE, "text/plain; charset=UTF-8");

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
	    }
	
	private Session createSession(Request httpRequst, Response httpResponse){
		
		String sessionId = null;
		boolean found = false;
		
		//check for already existing cookie
		//cookie format -> Cookie: puppy-io.sessionid=XXX; .....
		String cookieHeader = httpRequst.getHeader("Cookie");
		
		if(cookieHeader != null){
			String[] cookies = cookieHeader.split(";");
			
			for(String checkStr : cookies){
				String pattern = ".*\\b" + SessionStore.sessionIdKey + "\\b=(.*)";

				Pattern r = Pattern.compile(pattern);
				Matcher m = r.matcher(checkStr);
				
				if(m.find()){
					sessionId = m.group(1);
					
					if(sessionStore.checkSessionIdExists(sessionId))
						found = true;
					
					break;
				}
				
			}
			
		}
		
		
		if(clusterConnector == null){//no cluster set
			
			if(!found){
				//add new user to session store
				sessionId = sessionStore.addNewUser();
				
				String cookie = SessionStore.sessionIdKey + "=" + sessionId;
				
				//add Cookie headers to response
				httpResponse.setHeader("Set-Cookie", cookie);
			}
			
			Session session = new SessionImpl(sessionStore, sessionId);
			return session;
		}else{
			if(!found){
				//add new user to session store
				sessionId = sessionStore.addNewUser();
				
				String cookie = SessionStore.sessionIdKey + "=" + sessionId;
				
				//add Cookie headers to response
				httpResponse.setHeader("Set-Cookie", cookie);
			}
			
			Session session = new ClusterSessionImpl(sessionStore, sessionId, clusterConnector.getMulticastGroup());
			return session;
		}
		
		
	}
	
	
	public void requestProcess(ServerContext serverContext, String newLocation, com.lovi.quebic.web.enums.HttpMethod httpMethod) throws Exception {
		
		if (requestMapper != null) {
			String userRequestPath = null;
			boolean found = false;
			for (RequestMap requestMap : requestMapper.getRequestMaps()) {

				String requestPath = requestMap.getPath();
				String regExpPath = requestMap.getRegExpPath();

				String incomePath = newLocation.split("\\?")[0];// remove
																				// query
																				// parameters;
				incomePath = (incomePath.charAt(incomePath.length() - 1) == '/')? incomePath.substring(0, incomePath.length() - 1) : incomePath;// remove
																						// last
																						// slash

				userRequestPath = incomePath;
				
				// add index point
				requestPath = "/_INDEX_" + requestPath;
				regExpPath = "/_INDEX_" + regExpPath;
				incomePath = "/_INDEX_" + incomePath;

				String[] splitRequestPath = requestPath.split("/");
				String[] splitImcomePathStr = incomePath.split("/");

				// http method are not equal
				if (!(requestMap.getHttpMethod() == httpMethod))
					continue;

				// check route and incoming request path separators not equal
				if ((splitRequestPath.length != splitImcomePathStr.length))
					continue;

				// match
				Matcher matcher = Pattern.compile(regExpPath).matcher(incomePath);
				if (matcher.matches()) {

					// extract path variables
					for (int i = 0; i < splitRequestPath.length; i++) {

						String str = splitRequestPath[i];

						Matcher matcherParm = Pattern.compile("\\{(.*)\\}").matcher(str);

						if (matcherParm.matches()) {

							try {
								serverContext.getHttpRequst().getParameters().put(matcherParm.group(1),
										splitImcomePathStr[i]);
							} catch (Exception e) {
								continue;
							}
						}

					}

					//handler process
					Future<?> future = Future.create();
					try {
						requestMap.getHandler().handle(serverContext, future);
						
						if(future.isFail())
							throw future.getFailure();
						
					} catch (Throwable e) {
						
						int responseCode = 500;
						if(future.getErrorCode() != null)
							responseCode = future.getErrorCode();
						
						FailureHandler<ServerContext> failureHandlerRequestMap = requestMap.getFailureHandler();

						if (failureHandlerRequestMap != null)
							failureHandlerRequestMap.handle(serverContext, e, responseCode,
									ErrorMessage.INTERNAL_SERVER_ERROR.getMessage());
						else {
							FailureHandler<ServerContext> failureHandlerRequestMapper = requestMapper.getFailureHandler();
							if (failureHandlerRequestMapper != null)
								failureHandlerRequestMapper.handle(serverContext, e, responseCode,
										ErrorMessage.INTERNAL_SERVER_ERROR.getMessage());
							else
								writeServerError(serverContext, e, responseCode,
										ErrorMessage.INTERNAL_SERVER_ERROR.getMessage());
						}
					}
					found = true;
					break;
					
				}
			}

			if (!found) {
				FailureHandler<ServerContext> failureHandlerRequestMapper = requestMapper.getFailureHandler();
				if (failureHandlerRequestMapper != null)
					failureHandlerRequestMapper.handle(
							serverContext, 
							new ResourceNotFoundException(ErrorMessage.RESOURCE_NOT_FOUND.getMessage() + " " + userRequestPath), 
							404,
							ErrorMessage.RESOURCE_NOT_FOUND.getMessage());
				else
					writeServerError(
							serverContext, 
							new ResourceNotFoundException(ErrorMessage.RESOURCE_NOT_FOUND.getMessage() + " " + userRequestPath), 
							404,
							ErrorMessage.RESOURCE_NOT_FOUND.getMessage());
			}

		} else {
			throw new RequestMapperException();
		}

	}
	
	public void requestProcess(ServerContext serverContext, String newLocation) throws Exception {
		requestProcess(serverContext, newLocation, com.lovi.quebic.web.enums.HttpMethod.GET);
	}

}
