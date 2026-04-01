FROM golang:1.24 AS builder

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/saiao ./cmd/saiao

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/saiao /app/saiao

EXPOSE 8080
ENTRYPOINT ["/app/saiao"]
