package com.lovi.quebic.annotation;

import java.lang.annotation.Documented;
import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

import com.lovi.quebic.annotation.enums.ParmsType;

/**
 * Indicates that an annotated method is a "Service Function"
 * <br/>
 * <br/>
 * value => name of the Service Function. If you didn't provide a name for the Service Function then name of the annotated method is become the Service Function name. 
 * <br/>
 * <br/>
 * @author Tharanga Thennakoon
 *
 */

@Documented
@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
public @interface ServiceFunction {
	String value() default "";
	ParmsType inputParm() default ParmsType.PRIMITIVE;
}
