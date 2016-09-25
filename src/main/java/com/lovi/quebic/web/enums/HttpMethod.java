package com.lovi.quebic.web.enums;

public enum HttpMethod {
	GET, POST, PUT, DELETE, PATCH, HEAD;

	public static HttpMethod lookup(String method) {
		if (method == null)
			return null;

		try {
			return valueOf(method);
		} catch (IllegalArgumentException e) {
			return null;
		}
	}
}
