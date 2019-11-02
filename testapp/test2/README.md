## GO build
* CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .

## testapp docker build
* sudo docker build --no-cache -t quebicdocker/testapp:1.0.0 .
* sudo docker login
* sudo docker push quebicdocker/testapp:1.0.0

## run on k8
* sudo kubectl apply -f deployment.yaml
* sudo kubectl expose deployment testapp --target-port=3000 --type=NodePort

## get
* sudo kubectl get svc
* sudo kubectl get deployment

## scale down
* sudo kubectl scale deployment testapp --replicas=0

## delete
* sudo kubectl delete svc testapp
* sudo kubectl delete deployment testapp