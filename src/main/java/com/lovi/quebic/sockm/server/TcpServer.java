package com.lovi.quebic.sockm.server;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.channels.ServerSocketChannel;
import java.nio.channels.SocketChannel;
import java.util.Iterator;
import java.util.Set;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.sockm.config.TcpAddress;
import com.lovi.quebic.sockm.exception.TcpServerException;
import com.lovi.quebic.sockm.launcher.SockMLauncher;

public class TcpServer {

	private final Logger logger = LoggerFactory.getLogger(TcpServer.class);

	private TcpAddress localTcpAddress;
	private SockMLauncher launcher;
	
	private Thread serverThread;

	public TcpServer(SockMLauncher launcher) {
		this.launcher = launcher;
		
		localTcpAddress = this.launcher.getTcpServerAddress();
	}
	
	public void listen() throws TcpServerException {
		
		final ServerSocketChannel serverSocket;
		final Selector selector;
		
		try{

			if(localTcpAddress == null){
				
				serverSocket = ServerSocketChannel.open();
				serverSocket.bind(new InetSocketAddress(0));
				serverSocket.configureBlocking(false);
				
				selector = Selector.open();
				serverSocket.register(selector, SelectionKey.OP_ACCEPT);
				
				String address = serverSocket.socket().getInetAddress().getHostAddress();
				int port = serverSocket.socket().getLocalPort();
				localTcpAddress = new TcpAddress(address, port);
			}else{
				
				serverSocket = ServerSocketChannel.open();
				serverSocket.bind(new InetSocketAddress(localTcpAddress.getAddress(), localTcpAddress.getPort()));
				serverSocket.configureBlocking(false);
				
				selector = Selector.open();
				serverSocket.register(selector, SelectionKey.OP_ACCEPT);
				
			}
			
			logger.info("TCP Server Running : " + localTcpAddress);
			
		}catch(Exception e){
			throw new TcpServerException(e.getMessage());
		}
		
		serverThread = new Thread(new Runnable() {

			@Override
			public void run() {
				
				try{
					
					while (true) {
						
						//check for stop
						if(Thread.interrupted())
							break;

						selector.select();

						Set<SelectionKey> selectionKeys = selector.selectedKeys();
						Iterator<SelectionKey> iterator = selectionKeys.iterator();

						while (iterator.hasNext()) {
							SelectionKey key = iterator.next();
							iterator.remove();
							
							if(key.isValid()){
								
								if (key.isAcceptable()) {
									SocketChannel clientSocketChannel = serverSocket.accept();

									clientSocketChannel.configureBlocking(false);

									clientSocketChannel.register(selector, SelectionKey.OP_READ);

								} else if (key.isReadable()) {
									
									try{
										processTcpMessage(key);
									}catch(Exception e){
										
									}
									
								} else if (key.isWritable()) {
									SocketChannel clientSocketChannel = null;
									try{
										clientSocketChannel = (SocketChannel) key.channel();
										
										ByteBuffer buffer = ByteBuffer.allocateDirect(1024);
										if(!buffer.hasRemaining()){
											buffer.compact();
											key.interestOps(SelectionKey.OP_READ);
										}
										
										clientSocketChannel.close();
									}catch(Exception e){
										clientSocketChannel.close();
									}
								}
							}
						}
						
					}
				} catch (Exception e) {
					logger.error("tcpListen Process " + e.getMessage());
				}
			}
			
			
		});
		
		serverThread.start();
		
	}
	
	private void processTcpMessage(SelectionKey key) throws IOException{
		SocketChannel clientSocketChannel = (SocketChannel) key.channel();
		
		ByteBuffer clientBuffer = ByteBuffer.allocateDirect(1024);

		StringBuilder requestStringBuilder = new StringBuilder();

		int bytesRead = clientSocketChannel.read(clientBuffer); // read into buffer.
		
		if(bytesRead == -1){
			clientSocketChannel.close();
		}
		
		while (bytesRead > 0) {

			clientBuffer.flip();

			while (clientBuffer.hasRemaining()) {
				requestStringBuilder.append((char) clientBuffer.get());
			}

			clientBuffer.clear();

			bytesRead = clientSocketChannel.read(clientBuffer);

		}
		
		System.out.println("Input " + requestStringBuilder.toString());
		
		clientSocketChannel.write(ByteBuffer.wrap("I am server".getBytes()));
		key.interestOps(SelectionKey.OP_WRITE);
	}

	public TcpAddress getLocalTcpAddress() {
		return localTcpAddress;
	}

	
	
	public synchronized void stopServer(){
		if(serverThread != null){
			serverThread.interrupt();
		}
	}
	
}
