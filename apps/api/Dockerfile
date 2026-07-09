FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /bin/api /app/api
COPY scripts/migrations /app/scripts/migrations
ENV HTTP_ADDR=:8080
ENV SEED=1
WORKDIR /app
EXPOSE 8080
CMD ["/app/api"]
