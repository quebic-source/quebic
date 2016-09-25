# quebic
reactive java web-framework

quebic contains non-blocking web server. You can develop micro services by using quebic and you can communicate with each services easily. quebic have inbuilt Clustering and Load-Balancing mechanism.

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
	private ServiceCaller serviceCaller;
	
	@ResponseBody
	@RequestMapping(produce="application/json")
	public void findAll(HttpResponseResult responseResult) throws ServiceCallerException{
		
		Result<List<User>> result = Result.create();
		FailResult failResult = FailResult.create();
		
		serviceCaller.call("UserService.findAll", result);
		
		result.process(r->{
			responseResult.complete(new ResponseMessage(1, r));
		}, failResult);
		
		failResult.setHandler(fail->{
			responseResult.complete(new ResponseMessage(-1, fail.getMessage()),500);
		});
	}
	
	@ResponseBody
	@RequestMapping(method=HttpMethod.POST, produce="application/json")
	public void insert(@ModelAttribute User user, HttpResponseResult responseResult) throws ServiceCallerException{
		
		serviceCaller.call("UserService.insert", user);
		responseResult.complete(new ResponseMessage(1, "do insert"),200);
	
	}
	....
}
```
####@Controller
* Use ```com.lovi.puppy.annotation.Controller```
* Implementation of the controllers are similar to the spring-mvc but remember internal architecture of the puppy-io is totally different from spring-mvc

####@RequestMapping
* value = The primary mapping expressed by this annotation
* method = The HTTP request methods
* consumes = The consumable media types of the mapped request
* produce = The producible media types of the mapped request

####HttpResponseResult
* ```HttpResponseResult.complete(Object value)``` set response value
* ```HttpResponseResult.complete(Object value, int statusCode)``` set response value with statusCode
* If you put ```@ResponseBody``` annonation with the method, then return the value of object as response. otherwise response is redirect to  a template or another route.
* ```HttpResponseResult.complete("{template}")```
* ```HttpResponseResult.complete("/{route}")```
* puppy-io use Thymeleaf template engine for genarating templates

####ServiceCaller
* ServiceCaller is used to call service method
* ```ServiceCaller.call(String serviceMethod, Object... inputParameters)```
* ```ServiceCaller.call(String serviceMethod, Result<U> result, Object... inputParameters)```. if your service method has a return value, get the return value by using ```Result```
* ```ServiceCaller.call(String appName, String serviceMethod, Result<U> result, Object... inputParameters)```. if you want to call serivce method from a another application, you can call with the appName

####Result
* ```Result<T>``` is used to catch the return value of the service method which is called by ```ServiceCaller```

####FailResult 
* ```FailResult``` is used to catch the failure within the service method call which is called by ```ServiceCaller```

####ViewAttribute
```java
@Controller
public class IndexController {

	@Autowired
	private ServiceCaller serviceCaller;
	
	@RequestMapping
	public void loadIndexView(Session session, ViewAttribute viewAttribute, HttpResponseResult responseResult){
		User loggedUser = session.get("user", User.class);
		if(loggedUser != null){
			viewAttribute.put("loggedUser", loggedUser);
			responseResult.complete("users-dashboard");
		}else
			responseResult.complete("index");
	}
	....
}
```
* Use ```com.lovi.puppy.web.ViewAttribute```
* ```viewAttribute.put("loggedUser", loggedUser);```
* ```viewAttribute.get("loggedUser", User.class);```
* In the template ```${loggedUser.userId}```
* ViewAttribute is used to maintain any data that you want to share between handlers or share between views

####Session
```java
@Controller
public class IndexController {

	@RequestMapping(method=HttpMethod.POST)
	public void signIn(@RequestParm("userId") String userId, @RequestParm("password") String password,
					Session session,
						HttpResponseResult responseResult) throws ServiceCallerException{
		
		Result<User> result = Result.create();
		FailResult failResult = FailResult.create();
		
		serviceCaller.call("UserService.findByUserIdAndPassword", result, userId, password);
		
		result.process(user->{
			if(user != null)
				session.put("user", user);
			
			responseResult.complete("/");
			
		}, failResult);
		
		failResult.setHandler(fail->{
			responseResult.complete("/");
		});
		
	}
}
```
* Use ```com.lovi.puppy.web.Session```
* ```session.put("user", user);```
* ```User loggedUser = session.get("user", User.class);```

##Service App
```java
@Service("userService")
public class UserService{

	@Autowired
	private UserRepository userRepository;
	
	@ServiceFunction
	public void insert(User user){
		userRepository.insert(user);
	}
	
	@ServiceFunction("_findAll")
	public List<User> findAll(){
		return userRepository.findAll();
	}
	....
}
```
####@Service
* Use ```com.lovi.puppy.annotation.Service```
* Class is marked as a service by using ```@Service```

####@ServiceFunction
* Use ```com.lovi.puppy.annotation.ServiceFunction```
* Method is marked as a service method by using ```@ServiceFunction```
