FROM golang:1.10-stretch
MAINTAINER Radu Topala <radu.topala@trisoft.ro>

RUN apt-get update
RUN apt-get install -y unzip libgconf-2-4 chromium bzip2
RUN mv /usr/bin/chromium /usr/bin/chrome
VOLUME /code
ENV GOPATH /code
