package com.lovi.quebic.annotation;

import java.lang.annotation.Documented;
import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * Annotation which indicates that a method parameter should be bound to a web request parameter.
 * <br/>
 * <br/>
 * value => The name of the request parameter to bind to
 * <br/>
 * <br/>
 * defaultValue => The default value for request parameter to bind to
 * <br/>
 * <br/>
 * required => Whether the parameter is required.Defaults to true, leading to an exception being thrown if the parameter is missing in the request. Switch this to false if you prefer a null value if the parameter is not present in the request
 * <br/>
 * <br/>
 * @author Tharanga Thennakoon
 *
 */
@Documented
@Target(ElementType.PARAMETER)
@Retention(RetentionPolicy.RUNTIME)
public @interface RequestParm {
	String value();
	String defaultValue() default "";
	boolean required() default true;
}
