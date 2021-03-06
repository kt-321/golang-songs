FROM golang:latest as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/golang-songs
COPY . .
RUN make

# runtime image
FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/golang-songs/app /app
COPY --from=builder /go/src/golang-songs/db/ /db/
EXPOSE 8080
ENTRYPOINT ["/app"]
