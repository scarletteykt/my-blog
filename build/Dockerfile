FROM golang:1.18.3-alpine
ARG OUTPUT_BINARY
ARG BUILD_DIR
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
USER 0
ADD . /app/
WORKDIR /app
RUN apk update && apk add --no-cache bash git make
RUN make build
RUN mkdir -p /go
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN git clone https://github.com/vishnubob/wait-for-it.git
RUN ["chmod", "+x", "./build/entrypoint.sh"]
ENTRYPOINT ["./build/entrypoint.sh"]
CMD ["sh","-c","make run" ]