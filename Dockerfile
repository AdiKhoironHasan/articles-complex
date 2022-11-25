# copy image golang
FROM golang:1.17-alpine3.13 AS build-env

# create label
LABEL maintainer="adi khoiron hasan<adikhoironhasan@gmail.com>"

# create app name
ENV APP_NAME=articles-complex

# linux update
RUN apk update && apk upgrade && apk add git

RUN ls -ls

# make dicectory for service
RUN mkdir -p /src/articles-complex

# copy from computer to docker image
COPY . /src/articles-complex

# set work directory for now
WORKDIR /src/articles-complex

# get all dependemcy
RUN go get

# import the required dependencies and remove the unnecessary ones
RUN go mod tidy

# install golang migration
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# build image docker
RUN go build

# expose internal port service app
EXPOSE 8000

# run the service
CMD "./articles-complex"