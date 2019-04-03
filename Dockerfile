FROM golang:1.12

ADD ./app /watchdog/app
ADD ./internal /watchdog/internal
ADD ./go.mod /watchdog/go.mod
# ADD ./go.sum /authorization/go.sum

WORKDIR /watchdog

RUN go mod download
RUN go build -o watchdog_server ./app/server/main.go
RUN go build -o watchdog_watcher ./app/watcher/main.go
