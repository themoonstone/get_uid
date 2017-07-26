#FROM iron/go:dev
#MAINTAINER Razil "torytan@whispir.cc"
#WORKDIR /Users/sprixin/go/src/get_uid
#WORKDIR $GOPATH/src/get_uid
#ADD . $GOPATH/src/get_uid
#ADD . /Users/sprixin/go/src/get_uid
#RUN curl https://glide.sh/get | sh
#RUN go build .
#EXPOSE 8080
#ENTRYPOINT ["./get_uid"]


FROM iron/go

WORKDIR /get_uid

ADD . /get_uid

EXPOSE 8080

ENTRYPOINT ["./get_uid"]