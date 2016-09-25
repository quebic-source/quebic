package com.lovi.quebic.annotation;

import java.lang.annotation.Documented;
import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

import com.lovi.quebic.annotation.enums.ResponseBodyFormat;

/**
 * Annotation that indicates a method return value should be bound to the web response body.
 * @author Tharanga Thennakoon
 *
 */

@Documented
@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
public @interface ResponseBody {
	ResponseBodyFormat value() default ResponseBodyFormat.JSON;
}
