# Quebic - FaaS Framework

Quebic is a framework for writing serverless functions to run on Dockers or Kubernetes. You can write your functions in any language. Currently quebic supports only for Java and NodeJS.

![quebic](https://github.com/quebic-source/quebic/blob/master/docs/quebic.png)

### Getting Started

#### Install Docker
 * If you have already setup docker on your envirnment ,skip this step.
 * [Install Docker](https://docs.docker.com/install/)

#### Getting Binaries

###### For Linux Users
 * Download binaries from [here](https://github.com/quebic-source/quebic/blob/master/bin/quebic.tar.gz). Save and extract it into preferred location.
 * After extract, you can see quebic-mgr and quebic cli inside that dir. 
 
###### For Windows Users
 * [Install golang into your envirnment](https://golang.org/doc/install). 
 * Get [govendor](https://github.com/kardianos/govendor) tool. 
 * Run **govendor fetch**. This will download all the required dependencies for quebic.
 * Run for build quebic-mgr **go install quebic-faas/quebic-faas-mgr**
 * Run for build quebic cli **go install quebic-faas/quebic-faas-cli**
 * Congrats !!! Now you can find your binaries from $GOPATH/bin dir.

#### Run quebic-manager
 * Jump into quebic binaries location. Then run this commond **quebic-mgr**
 * By default quebic-mgr deploy its components ( eventbus, apigateway, functions ) as docker services. when you are in docker swrm manager it can deploy services which are created by quebic, among its cluster.
 * If you  want to deploy qubic  into kubernetes, set --deployment argumnet into kubernetes. 
 * Eg: **quebic-mgr --deployment kubernetes**
 * We will discuss more details about configurations in a later section. 
 
 
 
### Functions
#### Java Runtime
##### Programming Model
###### RequestHandler<Request, Response>
 * RequestHandler is an interface which comes with [quebic-runtime-java library](https://github.com/quebic-source/quebic-runtime-java). You can add your logic inside it's handle() method.
 * The Request Type and Response Type can be any Primitive datatype or Object.
```java
public class HelloFunction implements RequestHandler<Request, Response>{
 
 public void handle(Request request, CallBack<Response> callback, Context context) {
	 callback.success(new Response(1, "reply success"));
 }
 
}
```

###### Context
 * Context have these methods.
```java
BaseEvent baseEvent(); // return event details comes into this function
Messenger messenger(); // return messenger instance
Logger logger(); // return logger instance
```

###### CallBack
* CallBack provides way to reply.
```java
callback.success(); //reply 200 status code with empty reply data
callback.success("reply success"); //reply 200 status code with reply data
callback.success(201, "reply success"); //reply 201 status code with reply data

callBack.failure("Error occurred"); //reply 500 status code with reply err-data
callBack.failure(401, "Error occurred"); //reply 401 status code with reply err-data
```

###### Messenger
* Messenger provides way to publish events.
* void publish(String eventID, Object eventPayload, MessageHandler successHandler,ErrorHandler errorHandler,int timeout)throws MessengerException;
```java
context.messenger().publish("users.UserValidate", user, s->{
  user.setId(UUID.randomUUID().toString());
  callBack.success(201, user);
}, e->{
  callBack.failure(e.statuscode(), e.error());
}, 1000 * 8);

```
###### Logger
* Logger provides way to attach logs for perticular request context. We will discuss more about this logger in later section.
```java
context.logger().info("log info");
context.logger().error("log error");
context.logger().warn("log warn");
```

##### Create .jar artifact
 * Create new maven project.
 * Add this dependency and repository into .pom file.
 ```xml
<dependency>
    <groupId>com.quebic.faas.runtime</groupId>
    <artifactId>quebic-faas-runtime-java</artifactId>
    <version>0.0.1-SNAPSHOT</version>
</dependency>

<repositories>
    <repository>
     <id>quebic-runtime-java-mvn-repo</id>
     <url>https://raw.github.com/quebic-source/quebic-runtime-java/mvn-repo/</url>
     <snapshots>
      <enabled>true</enabled>
      <updatePolicy>always</updatePolicy>
     </snapshots>
    </repository>
</repositories>
```
 * Run **mvn clean package**
 
##### Deployment Spec
 * Create .yml spec file by describing how you want to deploy your functions into quebic. This is code snippet for deployment spec
 ```yml
function:
  name: hello-function # function name 
  artifactStoredLocation: /functions/hello-function.jar # jar artifact location
  handlerPath: com.quebicfaas.examples.HelloFunction # request handler java class
  runtime: java # function runtime
  events: # function going to listen these events
    - users.UserCreate
    - users.UserUpdate

route:
  requestMethod: POST
  url: /users
  requestMapping:
    - eventAttribute: eID
      requestAttribute: id
    ...
 ```
#### NodeJS Runtime

#### Manage your functions with quebic cli
##### Create function
* quebic function create --file [deployment spec file]
	
##### Update function
* quebic function update --file [deployment spec file]

##### Stop function
* quebic function stop --name [function name]
	
##### List all functions
* quebic function ls

##### Inspect function details
* quebic function inspect --name [function name]
	
