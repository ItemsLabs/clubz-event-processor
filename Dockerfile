FROM golang:1.20

LABEL maintainer="Sergey Shcherbina <trytoheal@gmail.com>"

WORKDIR $GOPATH/src/github.com/gameon-app-inc/fanclash-event-processor

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

CMD ["ufl_event_processor"]
