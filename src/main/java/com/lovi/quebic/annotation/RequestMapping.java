package com.lovi.quebic.annotation;

import java.lang.annotation.Documented;
import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

import com.lovi.quebic.web.enums.HttpMethod;

/**
 * Annotation for mapping web requests onto specific handler classes and/or handler methods
 * <br/>
 * <br/>
 * value => the path mapping URIs.(e.g. "/testPath" ,  "/testPath/:pathVariable"  )
 * <br/>
 * <br/>
 * method => The HTTP request methods to map (GET, POST, PUT, DELETE, PATCH). default GET. 
 * <br/>
 * <br/>
 * consumes => The consumable media types of the mapped request.
 * <br/>
 * <br/>
 * produce => The producible media types of the mapped request.
 * <br/>
 * <br/>
 * @author Tharanga Thennakoon
 *
 */

@Documented
@Target({ElementType.TYPE,ElementType.METHOD})
@Retention(RetentionPolicy.RUNTIME)
public @interface RequestMapping {

	String value() default "/";
	HttpMethod method() default HttpMethod.GET;
	String consumes() default "";
	String produce() default "";
	
}
