# quebic-faas-cli build
## linux
* GOOS=linux GOARCH=amd64 go build

## macos
* GOOS=darwin GOARCH=amd64 go build

## windows
* GOOS=windows GOARCH=amd64 go build

## build with script
### jump into quebic-faas-cli
* run ```sh build_script.sh <dist dir> <version>```
* run ```sh build_script.sh /home/dist/qb 0.1.5```