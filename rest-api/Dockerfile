FROM golang:1.23.1-alpine3.20

COPY main.go go.mod /app/
COPY info.cfg /etc/sample-api/
WORKDIR /app

RUN go build .
EXPOSE 8080


ENTRYPOINT [ "/app/rest-api" ]