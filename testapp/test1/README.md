# testapp build
#### go to testapp dir. then run following commond
##### http://blog.wrouesnel.com/articles/Totally%20static%20Go%20builds/
##### CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .

## testapp docker build
##### sudo docker build --no-cache -t quebicdocker/testapp:1.0.0-beta .
##### sudo docker login
##### sudo docker push quebicdocker/testapp:1.0.0-beta

## run on k8
### sudo kubectl run testapp --image=quebicdocker/testapp:1.0.0-beta --port=3000
### sudo kubectl expose deployment testapp --target-port=3000 --type=NodePort

## ingress 
### sudo kubectl apply -f node-ingress.yaml

## get
### sudo kubectl get svc
### sudo kubectl get deployment
### sudo kubectl get ingress

## delete
### sudo kubectl delete svc testapp
### sudo kubectl delete deployment testapp
### sudo kubectl delete ingress node-ingress