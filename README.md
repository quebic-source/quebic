# quebic
Faas framework

Quebic is a framework for writing serverless functions to run on Dockers or Kubernetes. You can write your functions in any language. Currently quebic supports for only Java and NodeJS.. [more](http://quebic.io/)

![quebic](https://github.com/quebic-source/quebic/blob/master/docs/quebic.png)

### Prerequisities
  * JDK 1.8.X
  * Maven 3.3.X

### Getting Started
 * Remote repository.
 
 ```xml
<repositories>
		<repository>
			<releases>
				<enabled>true</enabled>
				<updatePolicy>always</updatePolicy>
				<checksumPolicy>fail</checksumPolicy>
			</releases>
			<id>quebic_repo</id>
			<name>quebic_repo</name>
			<url>http://quebic.io/static/repo</url>
			<layout>default</layout>
		</repository>
	</repositories>
 ```
 
 * Add dependency.
 
```xml
	<dependency>
		<groupId>com.lovi.quebic</groupId>
		<artifactId>quebic-core</artifactId>
		<version>0.0.1-SNAPSHOT</version>
	</dependency>
```
 
### Sample application
 * Download the [sample-application](http://quebic.io/static/downloads/sample-app.zip)
 * Build the application using **mvn package**
 * Run the application using **java -jar target\sample-app-0.0.1-SNAPSHOT.jar**
 * Consume web app from **localhost:8080**

##Starting the application
```java
@SpringBootApplication
public class App 
{
	private final static Logger logger = LoggerFactory.getLogger(App.class);
    
    public static void main( String[] args )
    {
    	AppLauncher appLauncher = AppLauncher.create();
        
        appLauncher.run(App.class, 8080,r->{
        	logger.info(r);
        }, fail->{
        	logger.error(fail.getMessage());
        }, args);
        
    }
}
```

##Web App
```java
@Controller
@RequestMapping("/users")
public class UserController {

	@Autowired
	private UserService userService;
	
	@RequestMapping
	public void findAll(ApplicationContextData contextData, Response response) throws JsonParserException{
		response.writePOJO(userService.findAll(contextData));
	}
	
	@RequestMapping(value="/{id}")
	public void findById(@PathVariable("id") Integer id, ApplicationContextData contextData, Response response) throws JsonParserException{
		User user = userService.findById(contextData, id);
		response.writePOJO(user);
	}
	
	@RequestMapping(method=HttpMethod.POST)
	public void save(@ModelAttribute User user, ApplicationContextData contextData, Response response) throws JsonParserException{
		userService.save(contextData, user);
		response.setResponseCode(201);
		response.writePOJO(user);
	}
	
	....
}
```
####@Controller
* Use ```com.lovi.quebic.annotation.Controller```
* Implementation of the controllers are similar to the spring-mvc but remember internal architecture of the queubic is totally different from spring-mvc

####@RequestMapping
* value = The primary mapping expressed by this annotation
* method = The HTTP request methods
* consumes = The consumable media types of the mapped request
* produce = The producible media types of the mapped request

####Response
* ```Response.write(byte[] value)``` set response value
* ```Response.writPOJO(Object value)``` write java POJO as response.
