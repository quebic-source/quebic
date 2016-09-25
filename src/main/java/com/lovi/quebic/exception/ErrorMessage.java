package com.lovi.quebic.exception;

public enum ErrorMessage {

	//puppy-io core
	INTERNAL_SERVER_ERROR("INTERNAL_SERVER_ERROR"),
	UNABLE_TO_FOUND_REQUEST_MAPPER("UNABLE_TO_FOUND_REQUEST_MAPPER"),
	RESOURCE_NOT_FOUND("RESOURCE_NOT_FOUND"),
	REQUEST_MAP_MUST_START_WITH_SLASH("REQUEST_MAP_MUST_START_WITH_SLASH"), 
	REQUEST_FORWORD_ERROR("REQUEST_FORWORD_ERROR"),
	SESSION_USER_ALREADY_REMOVED("SESSION_USER_ALREADY_REMOVED"),
	SESSION_USER_NOT_FOUND("SESSION_USER_NOT_FOUND"),
	SESSION_PARM_NOT_FOUND("SESSION_PARM_NOT_FOUND => "),
	TEMPLATE_NAME_EMPTY("TEMPLATE_NAME_EMPTY"),
	TEMPLATE_NOT_FOUND("TEMPLATE_NOT_FOUND => "),
	STATIC_RESOURCE_NOT_FOUND("STATIC_RESOURCE_NOT_FOUND => "),
	EMPTY_REQUEST_ERROR("EMPTY_REQUEST_ERROR"),
	//puppy-io core
	
	REQUEST_MAPPING_ANNOTATION_VALUE_CAN_NOT_BE_EMPTY("@RequestMapping value can't be empty"),
	REQUEST_PARAM_ANNOTATION_VALUE_CAN_NOT_BE_EMPTY("@RequestParm value can't be empty"),
	REQUEST_PARAM_NOT_FOUND("request parameter not found -> "),
	UNABLE_TO_PARSE_REQUEST_PARM("unable to parse request parameter -> "), 
	PATH_PARAM_ANNOTATION_VALUE_CAN_NOT_BE_EMPTY("@PathVariable value can't be empty"), 
	PATH_PARAM_NOT_FOUND("path parameter not found"), 
	UNABLE_TO_FOUND_MODEL_ATTRIBUTE_CONSTRUCTOR ("unable to model attribute constructor -> "),
	UNABLE_TO_PARSE_PATH_PARM("unable to parse path parameter -> "),
	UNABLE_TO_PROCESS_METHOD_INPUT_PARM("unable to process method input parameter. method -> "),
	UNABLE_TO_PARSE_METHOD_RETURN_TYPE_VALUE("unable to parse method return type value. method -> "),
	UNABLE_TO_FOUND_SERVICE("No service for address -> "),
	SERVICE_CALL_UNABLE_TO_PROCESS("unable to process service call"),
	SERVICE_CALL_MESSAGE_CONTENT_NULL("message content is null for service call"),
	SERVICE_CALL_ILLEGAL_INPUT_PARAMETERS("illigel input parameters for service call"),
	SERVICE_CALL_RETURN_TYPE_MIS_MATCH("return type is not matched for service call"),
	SERVICE_ERROR("service fail -> "),
	UI_CALL_UNABLE_TO_PROCESS("unable to process ui call"),
	UI_SERVICE_FUNCTION_INPUT_PARAMETERS_FOUND("ui service function can't have input parameters"),
	UI_SERVICE_FUNCTION_UNABLE_TO_PROCESS("unable to process ui service function"),
	SESSION_OBJECT_UNABLE_TO_PARSE("unable to parse object into session. session key -> "), 
	SERVICE_CALL_UNABLE_TO_FOUND_DEFAULT_CONSTRUCTOR("unable to found default constructor of POJO"),
	SERVICE_CALL_UNABLE_TO_ACCESS_SERVICE_FUNCTION("unable to access service function"), 
	UNABLE_TO_REDIRECT_RESPONSE_TO_TEMPLATE("unable to redirect response to template"),
	UNABLE_TO_SAVE_FILE("unable to save file > "), 
	UNABLE_TO_FILE_STREAMS_COPY("UNABLE_TO_FILE_STREAMS_COPY ");
	
	String message;
	private ErrorMessage(String message){
		this.message = message;
	}
	
	public String getMessage(){
		return message;
	}
}
