FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -trimpath \
    -buildvcs=false \
    -ldflags="-s -w -buildid=" \
    -o /server .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /server /server

EXPOSE 3000

USER 65532:65532

ENTRYPOINT ["/server"]
