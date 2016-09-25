package com.lovi.quebic.annotation;

import java.lang.annotation.Documented;
import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * Indicates that an annotated method is a "UI Service Function"
 * <br/>
 * <br/>
 * value => Name of the UI Service Function. If you didn't provide a name for the UI Service Function then name of the annotated method is become the UI Service Function name. 
 * <br/>
 * <br/>
 * delay(second) => Indicates the UI Service Function firing rate in second.
 * <br/>
 * <br/>
 * @author Tharanga Thennakoon
 *
 */

@Documented
@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
public @interface UIServiceFunction {
	String value() default "";
	int delay();
}
