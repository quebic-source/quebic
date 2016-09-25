package com.lovi.quebic.servicecaller.impl;

import java.lang.annotation.Annotation;
import java.lang.reflect.Method;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.config.BeanDefinition;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.ClassPathScanningCandidateComponentProvider;
import org.springframework.core.type.filter.AnnotationTypeFilter;

import com.lovi.quebic.annotation.Service;
import com.lovi.quebic.annotation.ServiceFunction;
import com.lovi.quebic.annotation.enums.ParmsType;
import com.lovi.quebic.handlers.Handler;
import com.lovi.quebic.servicecaller.MessageBody;
import com.lovi.quebic.servicecaller.MessengerMapperGenerator;
import com.lovi.quebic.sockm.message.messenger.Message;
import com.lovi.quebic.sockm.messenger.Messenger;

public class MessengerMapperGeneratorImpl implements MessengerMapperGenerator {

	final static Logger logger = LoggerFactory.getLogger(MessengerMapperGenerator.class);
	
	private ApplicationContext applicationContext;
	
	public MessengerMapperGeneratorImpl() {
	}
	
	public MessengerMapperGeneratorImpl(ApplicationContext applicationContext) {
		this.applicationContext = applicationContext;
	}
	
	@Override
	public void start(Class<?> baseClass, Messenger messenger) throws Exception{
		ClassPathScanningCandidateComponentProvider scanner = new ClassPathScanningCandidateComponentProvider(false);
		scanner.addIncludeFilter(new AnnotationTypeFilter(Service.class));
		
		for(BeanDefinition bd : scanner.findCandidateComponents(baseClass.getPackage().getName())){
			
			Class<?> serviceAnnotatedClass = Class.forName(bd.getBeanClassName());
			
			Service serviceAnnotation = (Service) serviceAnnotatedClass.getAnnotation(Service.class);
			String serviceName = serviceAnnotation.value();
			if (serviceName.equals(""))
				serviceName = serviceAnnotatedClass.getSimpleName();
			
			final Object serviceObject;
			if(applicationContext == null)
				serviceObject = serviceAnnotatedClass.newInstance();
			else
				serviceObject = applicationContext.getBean(serviceAnnotatedClass);
			
			// process methods
			for(Method method : serviceAnnotatedClass.getDeclaredMethods()) {
				
				Annotation serviceFunctionMethodAnnotation = method.getAnnotation(ServiceFunction.class);

				// process @ServiceFunction
				if (serviceFunctionMethodAnnotation != null) {

					ServiceFunction serviceFunctionAnnotation = (ServiceFunction) serviceFunctionMethodAnnotation;

					String serviceFunctionName = serviceFunctionAnnotation.value();

					if (serviceFunctionName.equals(""))
						serviceFunctionName = method.getName();

					String serviceAddress = serviceName + "." + serviceFunctionName;

					
					//check inputParm type uding Annotation
					if(serviceFunctionAnnotation.inputParm() == ParmsType.PRIMITIVE){
						//default
						deployService(messenger, method, serviceAddress, serviceObject);
					}else{
						deployServicePojo(messenger, method, serviceAddress, serviceObject);
					}
					
					
				}
				
			}
			
		}
	}
	
	private void deployService(Messenger messenger, Method method, String serviceAddress, Object serviceObject) {
		
		messenger.subscribe(serviceAddress, MessageBody.class, Object.class, new Handler<Message<MessageBody,Object>>() {

			@Override
			public void handle(Message<MessageBody, Object> message) {
				
				try {
					Object[] consumeInputparameters = message.getMessage().getValues();
					Object[] inputParameters = new Object[method.getParameterCount()];

					for (int i = 0; i < inputParameters.length; i++) {
						try {
							inputParameters[i] = consumeInputparameters[i];
						} catch (Exception e) {
							throw new IllegalArgumentException();
						}
					}

					Object returnValue = method.invoke(serviceObject, inputParameters);
					message.reply(returnValue);

				} catch (Exception e) {
					message.replyFailure(e);
				}
			}
			
		});
		
		logger.info("service -> {}", serviceAddress);
		
		
	}
	
	private void deployServicePojo(Messenger messenger, Method method, String serviceAddress, Object serviceObject) {
		
		//Service POJO must have one input parameter.
		if(method.getParameterCount() != 1){
			logger.error("POJO based services must have only one input parameter");
			return;
		}
		
		Class<?> inputParmType = method.getParameterTypes()[0];
		Class<?> outputParmType = method.getReturnType();
		
		messenger.subscribe(serviceAddress, inputParmType, outputParmType, message->{
			
			try {
				Object consumeInputparameter = message.getMessage();
				Object inputParameter = consumeInputparameter;

				Object returnValue = method.invoke(serviceObject, inputParameter);
				message.reply(returnValue);

			} catch (Exception e) {
				message.replyFailure(e);
			}
			
		});
		logger.info("service -> {}", serviceAddress);
	}
}
