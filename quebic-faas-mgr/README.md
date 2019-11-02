# quebic-faas-mgr build
## quebic-faas-mgr go build
#### quebic-faas-mgr go build. jump into quebic-faas-mgr dir. then run following commond
##### http://blog.wrouesnel.com/articles/Totally%20static%20Go%20builds/
##### CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .

## quebic-faas-mgr docker build
##### sudo docker build --no-cache -t quebicdocker/quebic-faas-mgr:0.1.0 .
##### sudo docker login
##### sudo docker push quebicdocker/quebic-faas-mgr:0.1.0