FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

ARG CMD_PATH=service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/$CMD_PATH

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -S go-chujang && adduser -S go-chujang -G go-chujang
USER go-chujang
WORKDIR /home/go-chujang
COPY --from=builder /bin/app ./

EXPOSE 5000
ENTRYPOINT ["./app"]
