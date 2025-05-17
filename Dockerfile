FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /go/src/app

COPY . .

ENV CGO_ENABLED="0"

RUN go mod download
RUN go mod verify
RUN go build -o ping ./cmd/

FROM scratch AS final

LABEL authors="Sebastian Burmester"
LABEL org.opencontainers.image.source="https://github.com/sburmester/ping"

COPY --from=builder /go/src/app/ping ./

CMD ["./ping"]
