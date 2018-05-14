FROM alpine:3.4
MAINTAINER Justin K justin.j.kim94@gmail.com

RUN apk update
RUN apk add vim
RUN apk add curl



FROM mysql:5.7

FROM golang:1.8-onbuild


