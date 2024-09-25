FROM golang:1.23.1-alpine3.20

COPY . /app
WORKDIR /app
RUN mkdir -p /etc/sample-api/ && \
     echo "version 1.0.0" > /etc/sample-api/info.cfg

RUN go build .
EXPOSE 8080


ENTRYPOINT [ "/app/rest-api" ]