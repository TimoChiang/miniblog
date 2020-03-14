# Multiple Stage Builds
# First: In env with golang, create the program
FROM golang:alpine AS builder

#RUN apk update && apk upgrade && \
#

ADD . /go/src/miniblog
WORKDIR /go/src/miniblog
#RUN apk add git
# need the go.mod to download packages to local
RUN go mod download
RUN GOPATH=/go CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app miniblog

# Second: In small env, just run the program
FROM alpine:latest
COPY /views /blog/views
COPY --from=builder /go/src/miniblog/app /miniblog/app

# move to docker-compose.yaml
#ENV MYSQL_HOST "db"
# Need to move to blog to load the views folder
WORKDIR /miniblog
EXPOSE 8888
CMD ["/miniblog/app"]