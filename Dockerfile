FROM golang:alpine

WORKDIR /project/go-docker/

COPY go.* ./

COPY . .
RUN go build -o /project/go-docker/build/myapp .
#RUN mkdir -p /project/go-docker/build && go build -o /project/go-docker/build/myapp .

EXPOSE 8080

ENTRYPOINT [ "/project/go-docker/build/myapp" ]