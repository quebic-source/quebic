package com.lovi.quebic.cluster;

import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.channels.ServerSocketChannel;
import java.nio.channels.SocketChannel;
import java.util.Iterator;
import java.util.Set;
import java.util.concurrent.CompletableFuture;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.lovi.quebic.async.Future;
import com.lovi.quebic.cluster.loadblancer.LoadBlancer;
import com.lovi.quebic.common.Message;

class ProxyServer {

	final static Logger logger = LoggerFactory.getLogger(ProxyServer.class);
	
	private String sourceAddress;
	private int sourcePort;
	
	private static final int BUFSIZE = 1024;
	
	private LoadBlancer loadBlancer;
	
	ProxyServer(String sourceAddress,int sourcePort, LoadBlancer loadBlancer) {
		this.sourceAddress = sourceAddress;
		this.sourcePort = sourcePort;
		this.loadBlancer = loadBlancer;
	}
	
	void start(Future<String> future) throws Exception {
		
		ServerSocketChannel serverSocket = null;;
		Selector selector = null;
		try{
			serverSocket = ServerSocketChannel.open();
			serverSocket.bind(new InetSocketAddress(sourcePort));
			serverSocket.configureBlocking(false);
			
			selector = Selector.open();
			serverSocket.register(selector, SelectionKey.OP_ACCEPT);
			

			if (future != null)
				future.setResult(Message.MASTER_SERVER_START.getMessage() + sourceAddress + ":" + sourcePort);
			
		}catch(Exception e){
			throw e;
		}
		
		while(true){
			try{
				selector.select();
				Set<SelectionKey> keys = selector.selectedKeys();
				
				Iterator<SelectionKey> iterator = keys.iterator();
				while(iterator.hasNext()){
					
					SelectionKey key = iterator.next();
					iterator.remove();
					
					if(key.isValid()){
						
						if(key.isAcceptable()){
							accept(key);
						}
						
					}
				}
			}catch(Exception e){
				logger.error("start " + e.getMessage());
			}
		}
		
		
	}

	private void accept(SelectionKey selectionKey){
		SocketChannel clientRef = null;
		try{
			ServerSocketChannel serverSocketChannel = (ServerSocketChannel) selectionKey.channel();
			SocketChannel clientSocketChannel = serverSocketChannel.accept();
			clientRef = clientSocketChannel;
			
			clientSocketChannel.configureBlocking(false);
			
			clientSocketChannel.register(selectionKey.selector(), SelectionKey.OP_READ);
			
			Member destinationMember = loadBlancer.routeMember();//select destination member
			logger.info("select => \n" + destinationMember + "\n");
			
			InetSocketAddress destinationSocketAddress = new InetSocketAddress(destinationMember.getAddress(), destinationMember.getPort());
			SocketChannel destinationSocketChannel = SocketChannel.open(destinationSocketAddress);
			
			CompletableFuture.supplyAsync(() -> {
				forwardHandler(clientSocketChannel, destinationSocketChannel);
	    		return true;
			}).thenAccept(s -> {
				
			}).thenRun(() -> {}).exceptionally(fail->{
				logger.error("CompletableFuture exceptionally " + fail.getMessage());
				return null;
			});
		
		}catch(Exception e){
			logger.error("accept " + e.getMessage());
			
			try{
				clientRef.write(ByteBuffer.wrap(e.getMessage().getBytes()));
				clientRef.close();
			}catch(Exception ex){}
			
		}
		
	}
	
	private void forwardHandler(SocketChannel inChannel, SocketChannel outChannel){
		
        Selector sel = null;
        SocketChannel tmp;
        Set<SelectionKey> ready_keys;
        SelectionKey key;
        ByteBuffer transfer_buf = ByteBuffer.allocate(BUFSIZE);

        try {
            sel = Selector.open();
            
            inChannel.configureBlocking(false);
            outChannel.configureBlocking(false);
            
            inChannel.register(sel, SelectionKey.OP_READ);
            outChannel.register(sel, SelectionKey.OP_READ);
            
            while (true) {
            	
            	sel.select();
            	
                ready_keys=sel.selectedKeys();
                
                Iterator<SelectionKey> iterator = ready_keys.iterator();
                
                while(iterator.hasNext()) {
                    key  =(SelectionKey) iterator.next();
                    iterator.remove();
                    
                    tmp = (SocketChannel) key.channel();
                    if(tmp == null) {
                        continue;
                    }
                    
                    if (key.isReadable()) { 
                        if (tmp == inChannel) {
                           
                            if (relay(tmp, outChannel, transfer_buf) == false)
                                return;
                        }
                        if (tmp == outChannel) {
                         
                            if (relay(tmp, inChannel, transfer_buf) == false)
                                return;
                        }
                    }
                }
            }
        }
        catch (Exception ex) {
            logger.error("forwardHandler " + ex.getMessage());
        }
        finally {
            close(sel, inChannel, outChannel);
        }
		
	}
	
	private boolean relay(SocketChannel from, SocketChannel to, ByteBuffer buf){
        
		try{
        	int data = from.read(buf);
    		
    		if(data == -1)
    			return false;
    		
    		while(data > 0){
    			
    			buf.flip();

    			to.write(buf);
    			
    			buf.clear();
    			
    			data = from.read(buf);
    		}
    		
    		return true;
    		
        }catch(Exception e){
        	logger.error("relay " + e.getMessage());
        	return false;
        }
		
		
		
    }
	
	private void close(Selector sel, SocketChannel in_channel, SocketChannel out_channel) {
        try {
            if (sel != null)
                sel.close();
        }
        catch (Exception ex) {
        }
        try {
            if (in_channel !=null)
                in_channel.close();
        }
        catch (Exception ex) {
        }
        try {
            if (out_channel !=null)
                out_channel.close();
        }
        catch (Exception ex) {
        }
    }
	
}
