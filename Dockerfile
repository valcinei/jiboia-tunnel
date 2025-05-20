FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o jiboia-relay ./cmd/jiboia-relay

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/jiboia-relay /usr/bin/jiboia-relay

EXPOSE 80

ENTRYPOINT ["jiboia-relay"]
CMD ["--addr", ":80"]
