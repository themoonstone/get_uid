APPNAME=get_uid/get_uid

BUILD=docker run -it -v `pwd`:/go/src/$(APPNAME) -w /go/src/$(APPNAME) iron/go:1.7-dev go build

server:
	$(BUILD) -o get_uid
	
docker-build:server
	docker build -t get_uid .
	