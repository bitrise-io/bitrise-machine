FROM golang:1.7.3

ENV PROJ_NAME bitrise-machine

RUN apt-get update

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install git mercurial curl rsync ruby

#
# Install Bitrise CLI
RUN curl -fL https://github.com/bitrise-io/bitrise/releases/download/1.4.5/bitrise-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise
RUN chmod +x /usr/local/bin/bitrise
RUN bitrise setup

# Install required (testing) tools
#  Check for unhandled errors
RUN go get -u github.com/kisielk/errcheck
#  Go lint
RUN go get -u github.com/golang/lint/golint

# From the official Golang Dockerfile
#  https://github.com/docker-library/golang/blob/master/1.4/Dockerfile
RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p /go/src/github.com/bitrise-tools/$PROJ_NAME
COPY . /go/src/github.com/bitrise-tools/$PROJ_NAME

WORKDIR /go/src/github.com/bitrise-tools/$PROJ_NAME
# install
RUN go install

CMD $PROJ_NAME --version
