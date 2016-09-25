package com.lovi.quebic.annotation;

import java.lang.annotation.Documented;
import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * Indicates that an annotated class is a "UI Service"
 * <br/>
 * <br/>
 * value => name of the UI Service. this must be a unique value. If you didn't provide a name for the UI Service then name of the annotated class is become the UI Service name. 
 * <br/>
 * <br/>
 * @author Tharanga Thennakoon
 *
 */

@Documented
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
public @interface UIService {
	String value() default "";
}
