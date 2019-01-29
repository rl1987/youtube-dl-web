FROM golang:latest
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o main . 

FROM debian:latest
RUN apt-get -y update
RUN apt-get -y install youtube-dl
WORKDIR /root/
COPY --from=0 /go/src/app/main .
COPY --from=0 /go/src/app/public_html /root/public_html/
EXPOSE 8000
ENTRYPOINT /root/main
