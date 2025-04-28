FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

RUN apk update && apk upgrade
RUN apk add upx gcc musl-dev

WORKDIR /src
COPY . .

RUN go build -ldflags "-s -w" -o /bin/aac ./cmd/copy \
  && strip /bin/aac \
  && upx -q -9 /bin/aac
  
RUN go build -ldflags "-s -w" -o /bin/aac_cleanup ./cmd/cleanup \
  && strip /bin/aac_cleanup \
  && upx -q -9 /bin/aac_cleanup

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/aac /aac
COPY --from=builder /bin/aac_cleanup /aac_cleanup

ENTRYPOINT ["/aac"]