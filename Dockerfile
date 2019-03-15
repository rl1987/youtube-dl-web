FROM golang:latest
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o main . 

FROM debian:8
RUN apt-get -y update
RUN apt-get -y install curl python cron
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/bin/youtube-dl
RUN chmod a+rx /usr/bin/youtube-dl
ADD crontab /etc/cron.d/update-cron
RUN chmod 0644 /etc/cron.d/update-cron
RUN crontab /etc/cron.d/update-cron
CMD cron

WORKDIR /root/
COPY --from=0 /go/src/app/main .
COPY --from=0 /go/src/app/public_html /root/public_html/
EXPOSE 8000
HEALTHCHECK --interval=1m --timeout=1s CMD curl -f http://localhost:8000/ || exit 1
ENTRYPOINT /root/main
