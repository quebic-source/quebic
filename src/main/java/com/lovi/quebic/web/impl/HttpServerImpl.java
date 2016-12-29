package com.lovi.quebic.web.impl;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelOption;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslContextBuilder;
import io.netty.handler.ssl.util.SelfSignedCertificate;
import java.io.StringWriter;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.thymeleaf.TemplateEngine;
import org.thymeleaf.templateresolver.ClassLoaderTemplateResolver;

import com.lovi.quebic.async.Future;
import com.lovi.quebic.cluster.ClusterConnector;
import com.lovi.quebic.cluster.option.ClusterOption;
import com.lovi.quebic.cluster.option.MulticastGroup;
import com.lovi.quebic.common.Message;
import com.lovi.quebic.common.ServerThreadFactory;
import com.lovi.quebic.exception.ErrorMessage;
import com.lovi.quebic.exception.RequestMapperException;
import com.lovi.quebic.exception.TemplateException;
import com.lovi.quebic.handlers.FailureHandler;
import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.handlers.RequestHandler;
import com.lovi.quebic.web.ApplicationContextData;
import com.lovi.quebic.web.HttpServer;
import com.lovi.quebic.web.HttpServerOption;
import com.lovi.quebic.web.RequestMap;
import com.lovi.quebic.web.RequestMapper;
import com.lovi.quebic.web.Response;
import com.lovi.quebic.web.ServerContext;
import com.lovi.quebic.web.SessionStore;
import com.lovi.quebic.web.enums.HttpMethod;
import com.lovi.quebic.web.option.SSLServerConfig;
import com.lovi.quebic.web.server.HttpServerInitializer;
import com.lovi.quebic.web.template.TemplateOption;
import com.lovi.quebic.web.template.resource.StaticResourceOption;

/**
 * 
 * @author Tharanga Thennakoon
 *
 */
public class HttpServerImpl implements HttpServer {

	final static Logger logger = LoggerFactory.getLogger(HttpServer.class);
	private RequestMapper requestMapper;
	private SessionStore sessionStore;
	private ApplicationContextData applicationContextData;
	private List<HttpServerOption> httpServerOptions = new ArrayList<>();
	private ClusterConnector clusterConnector = null;
	private ClusterOption clusterOption;
	private SSLServerConfig sslServerConfig;
	private TemplateOption templateOption;
	private StaticResourceOption staticResourceOption;
	private String serverRunningHostName;
	private int serverRunningPort;
	
	private final String boss_thread_name = "quebic-http-boss-server";
	private final String worker_thread_name = "quebic-http-worker-server";
	private final int worker_server_thread_count = 100;
	
	public HttpServerImpl() {
		templateOption = TemplateOption.createThymeleafOption();
		staticResourceOption = StaticResourceOption.create();
	}
	
	@Override
	public void run(int port, RequestMapper requestMapper) {
		run("localhost", port, requestMapper, null, null);
	}

	@Override
	public void run(String hostname, int port, RequestMapper requestMapper) {
		run(hostname, port, requestMapper, null, null);
	}

	@Override
	public void run(int port, RequestMapper requestMapper, Handler<String> successHandler,
			Handler<Throwable> failureHandler) {
		run("localhost", port, requestMapper, successHandler, failureHandler);
	}

	@Override
	public void run(String hostname, int port, RequestMapper requestMapper, Handler<String> successHandler,
			Handler<Throwable> failureHandler) {
		
		Future<String> future = Future.create();
		future.setSussessHandler(successHandler);
		future.setFailureHandler(failureHandler);
		
		setRequestMapper(requestMapper);
		
		//check options
		if(httpServerOptions.size() > 0){
			for(HttpServerOption httpServerOption : httpServerOptions){
				
				if(httpServerOption instanceof ClusterOption){
					clusterOption = (ClusterOption) httpServerOption; 
				}else if(httpServerOption instanceof TemplateOption){
					templateOption = (TemplateOption) httpServerOption;
				}else if(httpServerOption instanceof StaticResourceOption){
					staticResourceOption = (StaticResourceOption) httpServerOption;
				}else if(httpServerOption instanceof SSLServerConfig){
					sslServerConfig = (SSLServerConfig) httpServerOption;
				}
				
			}
				
		}
		//check options end
		
		//process Cluster
		if(clusterOption != null){
			
			MulticastGroup multicastGroup = clusterOption.getMulticastGroup();
			if(multicastGroup == null)
				clusterConnector = new ClusterConnector();
			else
				clusterConnector = new ClusterConnector(multicastGroup.getMulticastAddress(), multicastGroup.getMulticastPort());
			
			if(clusterOption.isMaster()){
				
				try{
					
					if(clusterOption.getLoadBlancer() == null){
						throw new Exception("loadblancer not found");
					}
					
					logger.info(Message.MASTER_SERVER_CLUSTER_START.getMessage());
					serverRunningHostName = hostname;
					serverRunningPort = port;
					clusterConnector.startAsMaster(hostname, port, clusterOption.getLoadBlancer(), future);
				}catch(Exception e){
					if (future != null)
						future.setFailure(e);
				}
				
			}else{
				
				try{

					SslContext sslCtx = null;
					if(sslServerConfig != null){
						
						boolean SSL = System.getProperty("ssl") != null;
					    if(SSL){
					    	SelfSignedCertificate ssc = new SelfSignedCertificate();
					    	sslCtx = SslContextBuilder.forServer(ssc.certificate(), ssc.privateKey()).build();
					    } else {
					    	sslCtx = null;
					    }
						
					}
					
					EventLoopGroup bossGroup = new NioEventLoopGroup(1, new ServerThreadFactory(boss_thread_name));
					EventLoopGroup workerGroup = new NioEventLoopGroup(worker_server_thread_count, new ServerThreadFactory(worker_thread_name));
					
					ServerBootstrap b = new ServerBootstrap();
		            b.group(bossGroup, workerGroup)
		             .channel(NioServerSocketChannel.class)
		             .childHandler(new HttpServerInitializer(sslCtx, this))
		             .option(ChannelOption.SO_BACKLOG, 128)
		             .childOption(ChannelOption.SO_KEEPALIVE, true); 

		            prepareRequestMap();
		            
		            ChannelFuture f = b.bind(0).sync();
		            
		            serverRunningHostName = hostname;
		            
		            String localAddressStr = f.channel().localAddress().toString();
		            String[] localAddressStrSplit = localAddressStr.split(":");
		            
					serverRunningPort = Integer.parseInt(localAddressStrSplit[localAddressStrSplit.length - 1]);
					
					logger.info(Message.SERVER_CLUSTER_START.getMessage());
					
					applicationContextData = new ClusterApplicationContextDataImpl(clusterConnector.getMulticastGroup());
					sessionStore = new ClusterSessionStoreImpl(clusterConnector.getMulticastGroup());
					
					clusterConnector.setSessionStore(sessionStore);
					clusterConnector.setApplicationContextData(applicationContextData);
					
					clusterConnector.startAsWorker(serverRunningHostName, serverRunningPort);
					
		            if (future != null)
						future.setResult(Message.SERVER_START.getMessage() + serverRunningPort);
		            
		            f.channel().closeFuture().sync();
					
					
				}catch(Exception e){
					if (future != null){
						future.setFailure(e);
					}
				}
				
			}
			
		}else{
			serverRunningHostName = hostname;
			serverRunningPort = port;
			

			applicationContextData = new ApplicationContextDataImpl();
			sessionStore = new SessionStoreImpl();
			
			//no cluster
			setUpServer(future);
			
		}
		//process Cluster end

	}

	private void setUpServer(Future<String> future) {
		
		try {

			SslContext sslCtx = null;
			if(sslServerConfig != null){
				
				boolean SSL = System.getProperty("ssl") != null;
			    if(SSL){
			    	SelfSignedCertificate ssc = new SelfSignedCertificate();
			    	sslCtx = SslContextBuilder.forServer(ssc.certificate(), ssc.privateKey()).build();
			    } else {
			    	sslCtx = null;
			    }
				
			}
			
			EventLoopGroup bossGroup = new NioEventLoopGroup(1, new ServerThreadFactory(boss_thread_name));
			EventLoopGroup workerGroup = new NioEventLoopGroup(worker_server_thread_count, new ServerThreadFactory(worker_thread_name));
			
			ServerBootstrap b = new ServerBootstrap();
            b.group(bossGroup, workerGroup)
             .channel(NioServerSocketChannel.class)
             .childHandler(new HttpServerInitializer(sslCtx, this))
             .option(ChannelOption.SO_BACKLOG, 128)
             .childOption(ChannelOption.SO_KEEPALIVE, true); 

            prepareRequestMap();
            
            ChannelFuture f = b.bind(serverRunningPort).sync();
            
            if (future != null)
				future.setResult(Message.SERVER_START.getMessage() + serverRunningPort);
            
            f.channel().closeFuture().sync();
            
		} catch (Exception e) {
			if (future != null){
				future.setFailure(e);
			}
		}
	}

	private void prepareRequestMap() throws Exception {
		if (requestMapper != null) {
			
			prepareTemplateEngine();
			prepareStaticResources();
			
			for (RequestMap requestMap : requestMapper.getRequestMaps()) {

				StringBuilder regExpPath = new StringBuilder();

				String route = requestMap.getPath();

				if (route.charAt(0) != '/')
					throw new RequestMapperException(ErrorMessage.REQUEST_MAP_MUST_START_WITH_SLASH.getMessage());

				String[] splitRouteStr = route.split("/");

				// prepare reg ex
				for (int i = 0; i < splitRouteStr.length; i++) {

					String str = splitRouteStr[i];

					if (str.matches("\\{.*\\}"))
						regExpPath.append(".*");
					else if (str.matches("\\*")){
						regExpPath.append(".*");
						break;
					}
					else
						regExpPath.append(str);

					if ((i + 1) != splitRouteStr.length)
						regExpPath.append("/");
				}

				requestMap.setRegExpPath(regExpPath.toString());
			}
		} else
			throw new RequestMapperException();
	}
	
	
	@Override
	public void requestProcess(ServerContext serverContext, String newLocation) throws Exception {
		requestProcess(serverContext, newLocation, HttpMethod.GET);
	}

	@Override
	public void requestProcess(ServerContext serverContext, String newLocation, HttpMethod httpMethod) throws Exception {
		
	}
	
	@Override
	public void setRequestMapper(RequestMapper requestMapper) {
		this.requestMapper = requestMapper;
	}
	
	@Override
	public void addHttpServerOption(HttpServerOption option) {
		this.httpServerOptions.add(option);
	}

	private void prepareTemplateEngine(){
		
		logger.info("templates location -> " + templateOption.getLocation());
		
		//templates
		requestMapper.map(templateOption.getUrlAccess(), HttpMethod.GET).setHandler(new RequestHandler<ServerContext>() {
			
			@Override
			public void handle(ServerContext ctx, Future<?> future) {
				
				try{
					
					String templateName = (String) ctx.getData("templateName");
					ctx.removeData("templateName");
					
					if(templateName == null || templateName.equals(""))
						throw new TemplateException(ErrorMessage.TEMPLATE_NAME_EMPTY.getMessage());
					
					ClassLoaderTemplateResolver resolver = new ClassLoaderTemplateResolver();
		    		resolver.setTemplateMode(templateOption.getMode());
		    		resolver.setPrefix(templateOption.getLocation());
		    		resolver.setSuffix(templateOption.getSuffix());
		    		resolver.setCacheable(templateOption.getCacheable());
		    		
		    		
		    		TemplateEngine engine = new TemplateEngine();
		    		engine.setTemplateResolver(resolver);
		    		
		    		StringWriter writer = new StringWriter();
		    		
		    		try{
			    		engine.process(templateName, ctx.getTemplateContext(), writer);
		    		}catch (Exception e) {
		    			future.setErrorCode(404);
		    			throw new TemplateException(ErrorMessage.TEMPLATE_NOT_FOUND.getMessage() + templateName);
					}
					
		    		ctx.getHttpResponse().write(writer.toString());
		    		
				}catch(Exception e){
					future.setFailure(e);
				}
				
			}
		});
	}

	private void prepareStaticResources(){
		
		logger.info("static resource location -> " + staticResourceOption.getLocation());
		logger.info("static resource access url -> " + staticResourceOption.getUrlAccess() + "/*");
		
		//staticResources
		requestMapper.map(staticResourceOption.getUrlAccess() + "/*", HttpMethod.GET).setHandler(new RequestHandler<ServerContext>() {
			
			@Override
			public void handle(ServerContext ctx, Future<?> future) {
				String fileName = null;
				try{
					
					String incomeUrl = ctx.getHttpRequst().getLocation();
					fileName = incomeUrl.substring(staticResourceOption.getUrlAccess().length());
					
					Path file = Paths.get("src/main/resources/" + staticResourceOption.getLocation()  + fileName);
					
					byte[] bytes = Files.readAllBytes(file);
					
					String fileExt = fileName.substring(fileName.indexOf(".") + 1);
					
					Response response = ctx.getHttpResponse();
					
					if(fileExt.equals("js"))
						response.setContentType("application/javascript");
					else if(fileExt.equals("css"))
						response.setContentType("application/css");
					
					response.write(bytes);
		    		
				}catch(Exception e){
					String errorMsg = ErrorMessage.STATIC_RESOURCE_NOT_FOUND.getMessage() + staticResourceOption.getLocation()  + fileName;
					logger.error(errorMsg);
					future.setErrorCode(404);
					future.setFailure(new Throwable(errorMsg));
				}
				
			}
		}).setFailureHandler(new FailureHandler<ServerContext>() {
			
			@Override
			public void handle(ServerContext t, Throwable failure, int responseCode, String responseReason) {
				Response httpResponse = t.getHttpResponse();
				httpResponse.setResponseCode(responseCode);
				httpResponse.setResponseReason(failure.getMessage());
				httpResponse.write(failure.getMessage());
			}
		});
	}
	
	@Override
	public List<HttpServerOption> getHttpServerOption(){
		return httpServerOptions;
	}
	
	@Override
	public TemplateOption getTemplateOption(){
		return templateOption;
	}
	
	@Override
	public RequestMapper getRequestMapper(){
		return this.requestMapper;
	}
	
	@Override
	public ClusterConnector getClusterConnector(){
		return this.clusterConnector;
	}

	@Override
	public SessionStore getSessionStore() {
		return sessionStore;
	}

	@Override
	public ApplicationContextData getApplicationContextData() {
		return applicationContextData;
	}
	
	
}
