FROM golang
MAINTAINER yida
WORKDIR /go/src/
COPY . ./okex
EXPOSE 80
CMD ["/bin/bash", "/go/src/okex/script/build.sh"]