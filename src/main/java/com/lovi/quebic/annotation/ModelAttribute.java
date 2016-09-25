package com.lovi.quebic.annotation;

import java.lang.annotation.Documented;
import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;


/**
 * Annotation that binds a method parameter or method return value to a named model attribute, exposed to a web view
 * @author Tharanga Thennakoon
 *
 */
@Documented
@Target(ElementType.PARAMETER)
@Retention(RetentionPolicy.RUNTIME)
public @interface ModelAttribute {
	String value() default "";
}
