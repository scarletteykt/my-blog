FROM golang:1.18.3 as builder
ARG OUTPUT_BINARY
ARG BUILD_DIR
ADD . /app/
WORKDIR /app
RUN make build

FROM alpine:3.13.5
ARG OUTPUT_BINARY
USER 0
COPY --from=builder /app/$OUTPUT_BINARY /app/bin
COPY --from=builder /usr/local/go/ /usr/local/go/
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN git clone https://github.com/vishnubob/wait-for-it.git
RUN chown -R -f 888:888 /app && chmod +x /app
USER 888
CMD['/app/bin']