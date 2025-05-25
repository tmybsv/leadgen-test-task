FROM golang:1.24-alpine3.21 AS builder
RUN apk update && apk add --no-cache git make ca-certificates
RUN adduser -D -s /bin/sh hasher
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM alpine:3.21 AS compressor
RUN apk add --no-cache upx
COPY --from=builder /build/bin/hasher /bin/hasher
# RUN upx --best /bin/hasher
RUN upx /bin/hasher

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=compressor /bin/hasher /bin/hasher
COPY --from=builder /build/configs/ /configs/
USER hasher
EXPOSE 6969
ENTRYPOINT ["/bin/hasher"]
