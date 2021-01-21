FROM golang as builder

WORKDIR /go/src/github.com/BruceLEO1969/sqs-exporter
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sqs-exporter .

FROM alpine

COPY --from=builder /go/src/github.com/BruceLEO1969/sqs-exporter
RUN apk --update add --no-cache ca-certificates

EXPOSE 9434

CMD ["/sqs-exporter"]
