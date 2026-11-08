FROM golang:1.27.2-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /url-shortener .

FROM alpine:3.24
COPY --from=builder /url-shortener /usr/local/bin/url-shortener
EXPOSE 8080
CMD ["url-shortener"]
