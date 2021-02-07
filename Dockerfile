FROM golang

COPY . /snippetBox
WORKDIR /snippetBox

RUN go mod download
ENV GO111MODULE=on\
    CGO_ENABLED=0\
    GOOS=linux\
    GOARCH=amd64
RUN go build /snippetBox/cmd/web
RUN chmod +x /snippetBox/cmd/web
ENTRYPOINT ["./web"]

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /snippetBox .
CMD ["./web"]

EXPOSE 7070