# apigateway build
## apigateway go build
#### for apigateway go build. jump into quebic-faas-apigateway dir. then run following commond
##### http://blog.wrouesnel.com/articles/Totally%20static%20Go%20builds/
##### CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .

## apigateway docker build
##### sudo docker build --no-cache -t quebicdocker/quebic-faas-apigateway:1.0.0 .
##### sudo docker login
##### sudo docker push quebicdocker/quebic-faas-apigateway:1.0.0