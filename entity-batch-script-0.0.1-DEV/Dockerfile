FROM golang:1.12.6 as builder
ADD . /go/src/entity-batch-script
WORKDIR /go/src/entity-batch-script
RUN go get entity-batch-script
RUN go install
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
FROM scratch
COPY --from=builder /go/src/entity-batch-script/main /app/
WORKDIR /app
ENTRYPOINT ["./main"]