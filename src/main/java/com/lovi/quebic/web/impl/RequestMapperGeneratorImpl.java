package com.lovi.quebic.web.impl;

import java.lang.annotation.Annotation;
import java.lang.reflect.Method;
import java.lang.reflect.Parameter;
import java.util.Map;
import java.util.Map.Entry;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.config.BeanDefinition;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.ClassPathScanningCandidateComponentProvider;
import org.springframework.core.type.filter.AnnotationTypeFilter;
import org.springframework.expression.EvaluationContext;
import org.springframework.expression.Expression;
import org.springframework.expression.ExpressionParser;
import org.springframework.expression.spel.SpelParserConfiguration;
import org.springframework.expression.spel.standard.SpelExpressionParser;
import org.springframework.expression.spel.support.StandardEvaluationContext;

import com.lovi.quebic.annotation.Controller;
import com.lovi.quebic.annotation.ModelAttribute;
import com.lovi.quebic.annotation.PathVariable;
import com.lovi.quebic.annotation.RequestHeader;
import com.lovi.quebic.annotation.RequestMapping;
import com.lovi.quebic.annotation.RequestParm;
import com.lovi.quebic.async.Future;
import com.lovi.quebic.exception.ErrorMessage;
import com.lovi.quebic.exception.InternalServerException;
import com.lovi.quebic.exception.ModelAttibuteException;
import com.lovi.quebic.exception.RequestProcessingException;
import com.lovi.quebic.web.ApplicationContextData;
import com.lovi.quebic.web.Request;
import com.lovi.quebic.web.RequestMapper;
import com.lovi.quebic.web.RequestMapperGenerator;
import com.lovi.quebic.web.Response;
import com.lovi.quebic.web.ServerContext;
import com.lovi.quebic.web.Session;
import com.lovi.quebic.web.ViewAttribute;
import com.lovi.quebic.web.enums.HttpMethod;

public class RequestMapperGeneratorImpl implements RequestMapperGenerator{

	final static Logger logger = LoggerFactory.getLogger(RequestMapperGenerator.class);
	
	private ApplicationContext applicationContext;
	
	public RequestMapperGeneratorImpl() {
	}
	
	public RequestMapperGeneratorImpl(ApplicationContext applicationContext) {
		this.applicationContext = applicationContext;
	}
	
	@Override
	public void start(Class<?> baseClass, RequestMapper requestMapper) throws Exception{
		
		ClassPathScanningCandidateComponentProvider scanner = new ClassPathScanningCandidateComponentProvider(false);
		scanner.addIncludeFilter(new AnnotationTypeFilter(Controller.class));
		
		for(BeanDefinition bd : scanner.findCandidateComponents(baseClass.getPackage().getName())){
			
			Class<?> controllerAnnotatedClass = Class.forName(bd.getBeanClassName());
			
			String controllerBaseUrl = "";

			// processing annotations for class
			// process @RequestMapping
			Annotation requestMappingClassAnnonation = controllerAnnotatedClass.getAnnotation(RequestMapping.class);
			if (requestMappingClassAnnonation != null) {
				RequestMapping requestMappingAnnotation = (RequestMapping) requestMappingClassAnnonation;
				controllerBaseUrl = requestMappingAnnotation.value();
			}
			
			
			final Object controllerObject;
			if(applicationContext == null)
				controllerObject = controllerAnnotatedClass.newInstance();
			else
				controllerObject = applicationContext.getBean(controllerAnnotatedClass);
			
			// process methods
			for (Method method : controllerAnnotatedClass.getDeclaredMethods()) {

				Annotation requestMappingMethodAnnotation = method.getAnnotation(RequestMapping.class);
				
				// process @RequestMapping
				if (requestMappingMethodAnnotation != null) {
					RequestMapping requestMappingAnnotation = (RequestMapping) requestMappingMethodAnnotation;

					String requestUrl = controllerBaseUrl + requestMappingAnnotation.value();
					
					logger.info("web request mapping -> " + requestUrl + " | " + requestMappingAnnotation.method());

					switch (requestMappingAnnotation.method()) {
						case GET:
							requestMapper.map(requestUrl, HttpMethod.GET).setHandler((ctx, future) -> {
								processRequest(ctx, future, method, controllerObject, requestMappingAnnotation);
							});
							break;
						case POST:
							requestMapper.map(requestUrl, HttpMethod.POST).setHandler((ctx, future) -> {
								processRequest(ctx, future, method, controllerObject, requestMappingAnnotation);
							});
							break;
						case PUT:
							requestMapper.map(requestUrl, HttpMethod.PUT).setHandler((ctx, future) -> {
								processRequest(ctx, future, method, controllerObject, requestMappingAnnotation);
							});
							break;
						case DELETE:
							requestMapper.map(requestUrl, HttpMethod.DELETE).setHandler((ctx, future) -> {
								processRequest(ctx, future, method, controllerObject, requestMappingAnnotation);
							});
							break;
						default:
							requestMapper.map(requestUrl, HttpMethod.PATCH).setHandler((ctx, future) -> {
								processRequest(ctx, future, method, controllerObject, requestMappingAnnotation);
							});
							break;
					}
				}
			}
		}
	}

	private void processRequest(ServerContext serverContext, Future<?> future, Method method, Object controllerObject, RequestMapping requestMappingAnnotation){
		
		try{
			Request httpRequst = serverContext.getHttpRequst();
			Response httpResponse = serverContext.getHttpResponse();
			
			int methodtParameterCount = method.getParameterCount();
			Object[] inputParms = new Object[methodtParameterCount];
			int paramterCount = 0;
			
			// process through input parameters
			for(Parameter paramater : method.getParameters()) {
				
				String paramaterType = paramater.getType().getName();
				
				//check primitive type parameter
				if (paramaterType.equals(String.class.getName()) 
						|| paramaterType.equals(Integer.class.getName())
						|| paramaterType.equals(Double.class.getName())
						|| paramaterType.equals(Float.class.getName()) 
						|| paramaterType.equals(Long.class.getName())
						|| paramaterType.equals(Short.class.getName())
						|| paramaterType.equals(Boolean.class.getName())) {
					
					//@RequestParm
					RequestParm requestParm = paramater.getAnnotation(RequestParm.class);
					
					//@PathVariable
					PathVariable pathVariable = paramater.getAnnotation(PathVariable.class);
					
					//@RequestHeader
					RequestHeader requestHeader = paramater.getAnnotation(RequestHeader.class);
					
					//process @RequestParm
					if (requestParm != null) {
						String requestParmValue = requestParm.value();
						String requestParmDefaultValue = requestParm.defaultValue();
						boolean requestParmRequired = requestParm.required();
						
						if (requestParmValue.equals(""))
							throw new InternalServerException(ErrorMessage.REQUEST_PARAM_ANNOTATION_VALUE_CAN_NOT_BE_EMPTY.getMessage());
						else {
							String requestValue = httpRequst.getParameter(requestParmValue);

							if (requestValue == null){
								if(requestParmRequired && requestParmDefaultValue.equals(""))
									throw new RequestProcessingException(ErrorMessage.REQUEST_PARAM_NOT_FOUND.getMessage() + requestParmValue);
								else{
									if(requestParmDefaultValue.equals(""))
										requestValue = null;
									else
										requestValue = requestParmDefaultValue;
								}
							}
							try{
								if (paramaterType.equals(Integer.class.getName())) {
									inputParms[paramterCount++] = Integer.parseInt(requestValue);
								} else if (paramaterType.equals(Double.class.getName())) {
									inputParms[paramterCount++] = Double.parseDouble(requestValue);
								} else if (paramaterType.equals(Float.class.getName())) {
									inputParms[paramterCount++] = Float.parseFloat(requestValue);
								} else if (paramaterType.equals(Long.class.getName())) {
									inputParms[paramterCount++] = Long.parseLong(requestValue);
								} else if (paramaterType.equals(Short.class.getName())) {
									inputParms[paramterCount++] = Short.parseShort(requestValue);
								} else if (paramaterType.equals(Boolean.class.getName())) {
									inputParms[paramterCount++] = Boolean.parseBoolean(requestValue);
								} else {
									inputParms[paramterCount++] = requestValue;
								}
							}catch(Exception e){
								throw new RequestProcessingException(ErrorMessage.UNABLE_TO_PARSE_REQUEST_PARM.getMessage() + requestValue);
							}
							
						}
					}
					
					//process @PathVariable
					else if(pathVariable != null){
						String pathVariableValue = pathVariable.value();

						if (pathVariableValue.equals(""))
							throw new InternalServerException(ErrorMessage.PATH_PARAM_ANNOTATION_VALUE_CAN_NOT_BE_EMPTY.getMessage());
						else {
							String requestValue = httpRequst.getParameter(pathVariableValue);

							if (requestValue == null)
								throw new RequestProcessingException(ErrorMessage.PATH_PARAM_NOT_FOUND.getMessage() + pathVariableValue);

							try{
								if (paramaterType.equals(Integer.class.getName())) {
									inputParms[paramterCount++] = Integer.parseInt(requestValue);
								} else if (paramaterType.equals(Double.class.getName())) {
									inputParms[paramterCount++] = Double.parseDouble(requestValue);
								} else if (paramaterType.equals(Float.class.getName())) {
									inputParms[paramterCount++] = Float.parseFloat(requestValue);
								} else if (paramaterType.equals(Long.class.getName())) {
									inputParms[paramterCount++] = Long.parseLong(requestValue);
								} else if (paramaterType.equals(Short.class.getName())) {
									inputParms[paramterCount++] = Short.parseShort(requestValue);
								} else if (paramaterType.equals(Boolean.class.getName())) {
									inputParms[paramterCount++] = Boolean.parseBoolean(requestValue);
								} else {
									inputParms[paramterCount++] = requestValue;
								}
							}catch(Exception e){
								throw new RequestProcessingException(ErrorMessage.UNABLE_TO_PARSE_PATH_PARM.getMessage() + requestValue);
							}
							
						}
					}
					
					//process @RequestHeader
					else if(requestHeader != null){
						String requestHeaderValue = requestHeader.value();
						inputParms[paramterCount++] = httpRequst.getHeader(requestHeaderValue);
					}
					
					else{
						inputParms[paramterCount++] = null;
					}
					
				}
				
				else if (paramaterType.equals(ServerContext.class.getName())) {
					// check RoutingContext type parameter
					inputParms[paramterCount++] = serverContext;
				}

				else if (paramaterType.equals(Request.class.getName())) {
					// check HttpServerRequest type parameter
					inputParms[paramterCount++] = httpRequst;
				}
				
				else if (paramaterType.equals(Response.class.getName())) {
					// check HttpServerRequest type parameter
					inputParms[paramterCount++] = httpResponse;
				}
				
				else if (paramaterType.equals(Session.class.getName())) {
					// check Session type parameter
					Session session = serverContext.getSession();
					inputParms[paramterCount++] = session;
				}
				
				else if (paramaterType.equals(ApplicationContextData.class.getName())) {
					// check ApplicationContextData type parameter
					ApplicationContextData contextData = serverContext.getApplicationContextData();
					inputParms[paramterCount++] = contextData;
				}
				
				else if (paramaterType.equals(ViewAttribute.class.getName())) {
					// check ApplicationContextData type parameter
					ViewAttribute viewAttribute = new ViewAttributeImpl(serverContext);
					inputParms[paramterCount++] = viewAttribute;
				}
				
				else if (paramaterType.equals(Future.class.getName())) {
					// check Future type parameter
					inputParms[paramterCount++] = future;
				}
				
				else {
					// process @ModelAttribute
					ModelAttribute modelAttribute = paramater.getAnnotation(ModelAttribute.class);

					if (modelAttribute != null) {
						
						String modelAttributeValue = modelAttribute.value();

						Object parameterObject = null;
						try{
							parameterObject = paramater.getType().newInstance();
						}catch(IllegalAccessException | InstantiationException e){
							throw new ModelAttibuteException(ErrorMessage.UNABLE_TO_FOUND_MODEL_ATTRIBUTE_CONSTRUCTOR.getMessage() + e.getMessage());
						}

						// Turn on:
						// - auto null reference initialization
						// - auto collection growing
						SpelParserConfiguration config = new SpelParserConfiguration(true, true);
						ExpressionParser parser = new SpelExpressionParser(config);
						EvaluationContext context = new StandardEvaluationContext(parameterObject);

						if (modelAttributeValue.equals("")) {

							Map<String, String> parms = httpRequst.getParameters();
							for (Entry<String, String> entry : parms.entrySet()) {
								String key = entry.getKey();
								String value = entry.getValue();
								Expression exp = parser.parseExpression(key);
								
								try{
									exp.setValue(context, value);
								}catch(Exception e){
									//throw new Exception("ModelAttribute parse exception - unable to found parm - " + key);
								}
								
							}

						} else {

							Map<String, String> parms = httpRequst.getParameters();

							for (Entry<String, String> entry : parms.entrySet()) {

								// entry.getKey() => vehicle.owner.name
								// pattern => \\bvehicle\\.\\b
								// check for => vehicle.
								String checkStr = entry.getKey();
								String pattern = "\\b" + modelAttributeValue + "\\.\\b";

								Pattern r = Pattern.compile(pattern);

								Matcher m = r.matcher(checkStr);
								if (m.find()) {
									String key = checkStr.substring(m.end());
									String value = entry.getValue();
									Expression exp = parser.parseExpression(key);
									
									try{
										exp.setValue(context, value);
									}catch(Exception e){
										//throw new Exception("ModelAttribute parse exception - unable to found parm - " + key);
									}
								}
							}

						}

						inputParms[paramterCount++] = parameterObject;
					}else{
						throw new InternalServerException(ErrorMessage.UNABLE_TO_PROCESS_METHOD_INPUT_PARM.getMessage() + method.getName());
					}

				}
			}
			
			method.invoke(controllerObject, inputParms);
			
		}catch (ModelAttibuteException e) {
			future.setFail(400, e);
		}catch (Exception e) {
			future.setFailure(e);
		}
		
	}
	
	
	
}
