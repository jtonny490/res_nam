FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server
FROM node:22-alpine AS frontend-build
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend ./
RUN npm run build
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/server .
COPY --from=build /app/migrations ./migrations
COPY --from=frontend-build /app/frontend/dist ./frontend
COPY uploads ./uploads
EXPOSE 8080
CMD ["./server"]
