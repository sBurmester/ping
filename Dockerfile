FROM golang:1.24.2-alpine3.21 as builder

WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ping ./cmd

FROM scratch

LABEL authors="Sebastian Burmester"
LABEL org.opencontainers.image.source="https://github.com/sburmester/ping"

COPY --from=builder /go/src/app/ping ./

CMD ["./ping"]
