FROM golang:1.24.2-alpine3.21 as builder
LABEL authors="Sebastian Burmester"

WORKDIR /app
#COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ping .

FROM scratch
COPY --from=builder /app/ping .

EXPOSE 8080
CMD ["./ping"]
