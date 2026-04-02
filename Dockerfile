FROM golang:1.24 AS builder

WORKDIR /src
ARG VERSION=dev
ARG GIT_BRANCH=
ARG COMMIT=unknown
COPY . .
RUN effective_version="${VERSION}"; \
    if [ "${GIT_BRANCH}" != "master" ]; then effective_version="${VERSION}-dev"; fi; \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-X saiao/internal/buildinfo.Version=${effective_version} -X saiao/internal/buildinfo.Commit=${COMMIT}" \
    -o /out/saiao ./cmd/saiao

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/saiao /app/saiao

EXPOSE 8080
ENTRYPOINT ["/app/saiao"]
