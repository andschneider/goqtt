FROM golang:1.14.1 as build

WORKDIR /app
RUN go get github.com/andschneider/goqtt

COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine

WORKDIR /app
COPY --from=build /app/main /app/main

ENTRYPOINT ["/app/main"]
CMD ["-topic", "hello/world"]