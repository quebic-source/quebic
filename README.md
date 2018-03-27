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
#### Create Function
##### Java Runtime
###### Create .jar artifact
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
 
###### Deployment Spec
 * Create .yml spec file by describing how you want to deploy your functions into quebic.
 ```yml
function:
  name: hello-function
  artifactStoredLocation: /functions/quebic-faas-hellofunction-java/target/hello-function-0.0.1-SNAPSHOT.jar
  handlerPath: com.quebicfaas.examples.HelloFunction
  runtime: java
  events:
    - users.UserCreate

route:
  requestMethod: POST
  url: /users
  requestMapping:
    - eventAttribute: eID
      requestAttribute: id
    
    .....

 ```
    
 
#### RequestHandler<Request, Response>
 * RequestHandler is an interface which comes with quebic-runtime-java library. You can add your logic inside it's handle() method.

