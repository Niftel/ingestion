# Build Stage — compile on the native CI runner instead of emulating the target.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /praetor-ingestion .

# Run Stage
FROM alpine:3.24@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6

WORKDIR /

COPY --from=builder /praetor-ingestion /praetor-ingestion

RUN apk add --no-cache ca-certificates

EXPOSE 8081

USER 1000:1000

CMD ["/praetor-ingestion"]
