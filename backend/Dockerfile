FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.20
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
RUN adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /server ./server
COPY migrations ./migrations
RUN mkdir -p uploads && chown -R app:app /app
USER app
EXPOSE 8080
CMD ["./server"]
