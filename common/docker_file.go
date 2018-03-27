package common

//DockerFileContent_Java java docker file
const DockerFileContent_Java = "FROM quebicdocker/quebic-faas-container-java:1.0.0\nADD function.jar /app/function.jar\nARG access_key\nENV access_key $access_key\n"

//DockerFileContent_Java java docker file
const DockerFileContent_NodeJS = "FROM quebicdocker/quebic-faas-container-nodejs:1.0.0\nADD function_handler.tar /app/function_handler/\nARG access_key\nENV access_key $access_key\n"
